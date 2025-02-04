package repository

import (
	"context"
	"errors"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/aluko123/resume-go/internal/models"
)

//	var counter struct {
//		Seq int `bson:"seq"`
//	}
//var resumeIDCounter int64 = 0

// resume repository to handle all database ops for resumes
type ResumeRepository struct {
	collection *mongo.Collection
	//counterMutex sync.Mutex //mutex for synchronizing counter access
	nextResumeID int64 //in-memory copt of next ID
}

// create new resume repository
func NewResumeRepository(db *mongo.Database) *ResumeRepository {
	repo := &ResumeRepository{
		collection: db.Collection("resumes"),
	}

	//initialize nextResumeID from MongoDB
	err := repo.initializeNextResumeID()
	if err != nil {
		log.Fatal("Error initilaizing nextResumeID:", err)
	}

	return repo

}

// create goroutine to avoid race conditions
func (r *ResumeRepository) initializeNextResumeID() error {
	var counter struct {
		Seq int64 `bson:"seq"`
	}

	//r.counterMutex.Lock()
	//defer r.counterMutex.Unlock()

	err := r.collection.Database().Collection("counters").FindOne(
		context.Background(),
		bson.M{"_id": "resume_id"},
	).Decode(&counter)

	if err == nil { //counter document found
		r.nextResumeID = counter.Seq + 1 //initialize from DB
	} else if err == mongo.ErrNoDocuments { // counter document not found, so create it
		_, err = r.collection.Database().Collection("counters").InsertOne(
			context.Background(),
			bson.M{"_id": "resume_id", "seq": int64(0)},
		)

		if err != nil {
			return err
		}
		r.nextResumeID = 1 //start from 1
	} else {
		return err
	}
	return nil
}

// create new resume entry
func (r *ResumeRepository) Create(ctx context.Context, resume *models.Resume) error {
	//r.counterMutex.Lock()
	resume.ResumeID = int(atomic.AddInt64(&r.nextResumeID, 1) - 1) //assign the ID to the resume struct
	//nextID := r.nextResumeID              //copy the nextID before unlocking mutex
	//r.nextResumeID++
	//r.counterMutex.Unlock()

	//create this resume now
	resume.CreatedAt = time.Now()
	resume.UpdatedAt = time.Now()
	//log.Printf("ResumerIDCOunter: %d", resumeIDCounter)
	result, err := r.collection.InsertOne(ctx, resume)
	if err != nil {
		return err
	}

	resume.ID = result.InsertedID.(primitive.ObjectID) //set MongoDB ID

	//use goroutine to update func
	go func() {
		_, err := r.collection.Database().Collection("counters").UpdateOne(
			context.Background(),
			bson.M{"_id": "resume_id"},
			bson.M{"$set": bson.M{"seq": r.nextResumeID - 1}}, //decrement nextID
		)
		if err != nil {
			log.Println("Error updating counter in MongoDB:", err)
		}
	}()

	return nil
}

// get resume ID
func (r *ResumeRepository) GetByID(ctx context.Context, id string) (*models.Resume, error) {
	log.Printf("Resume id: %v", id)
	resumeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid resume ID: must be an integer")
	}

	var resume models.Resume
	err = r.collection.FindOne(ctx, bson.M{"resume_id": resumeID}).Decode(&resume)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("resume not found")
		}
		return nil, err
	}

	return &resume, nil
}

// update an existing resume
func (r *ResumeRepository) Update(ctx context.Context, resume *models.Resume) error {
	resume.UpdatedAt = time.Now()

	filter := bson.M{"_id": resume.ID}
	update := bson.M{"$set": resume}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("resume not found")
	}

	return nil
}

// delete removes a resume from database
func (r *ResumeRepository) Delete(ctx context.Context, id string) error {
	resumeID, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid resume ID: must be an integer")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"resume_id": resumeID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("resume not found")
	}

	return nil
}

// retreive all resumes
func (r *ResumeRepository) ListAll(ctx context.Context) ([]models.Resume, error) {
	//log.Println("ListAll repository method called") //debugging

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resumes []models.Resume
	if err = cursor.All(ctx, &resumes); err != nil {
		//log.Printf("Error decoding resumes: %v", err) //more debugging
		return nil, err
	}

	//log.Printf("Found resumes in DB: %v", resumes) //more debugging
	return resumes, nil
}

// update education entries
func (r *ResumeRepository) UpdateEducation(ctx context.Context, id string, education []models.Education) error {
	resumeID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"education":  education,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"resume_id": resumeID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("resume not found")
	}

	return nil
}

// update experience
func (r *ResumeRepository) UpdateExperience(ctx context.Context, id string, experience []models.Experience) error {
	resumeID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"experience": experience,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"resume_id": resumeID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("resume not found")
	}

	return nil
}

// update skills
func (r *ResumeRepository) UpdateSkills(ctx context.Context, resumeID string, skills models.Skills) error {
	objectID, err := primitive.ObjectIDFromHex(resumeID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"skills":     skills,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("resume not found")
	}

	return nil
}
