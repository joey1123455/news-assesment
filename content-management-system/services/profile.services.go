package services

import (
	"context"
	"errors"
	"time"

	"github.com/joey1123455/news-aggregator-service/content-management-system/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfileServices interface {
	FindUser(id string) (*models.UserProfile, error)
	CreateUser(id string, user *models.CreateUser) (*models.UserProfile, error)
	UpdateUser(id string, user *models.UpdateUser) (*models.UserProfile, error)
	DeleteUser(id string) error
}

type ProfileServicesImp struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewProfileService(ctx context.Context, collection *mongo.Collection) ProfileServices {
	return &ProfileServicesImp{
		ctx:        ctx,
		collection: collection,
	}
}

func (p *ProfileServicesImp) FindUser(id string) (*models.UserProfile, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.UserProfile

	query := bson.M{"_id": oid}
	err := p.collection.FindOne(p.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserProfile{}, err
		}
		return nil, err
	}

	return user, nil
}

func (p *ProfileServicesImp) CreateUser(id string, profile *models.CreateUser) (*models.UserProfile, error) {
	profile.ID, _ = primitive.ObjectIDFromHex(id)
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = profile.CreatedAt
	res, err := p.collection.InsertOne(p.ctx, profile)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("profile already exists")
		}
		return nil, err
	}

	var newPost *models.UserProfile
	query := bson.M{"_id": res.InsertedID}
	if err = p.collection.FindOne(p.ctx, query).Decode(&newPost); err != nil {
		return nil, err
	}

	return newPost, nil
}

func (p *ProfileServicesImp) UpdateUser(id string, user *models.UpdateUser) (*models.UserProfile, error) {

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: user}}
	res := p.collection.FindOneAndUpdate(p.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedProfile *models.UserProfile

	if err := res.Decode(&updatedProfile); err != nil {
		return nil, errors.New("no post with that Id exists")
	}

	return updatedProfile, nil
}

func (p *ProfileServicesImp) DeleteUser(id string) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": obId}

	res, err := p.collection.DeleteOne(p.ctx, query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}

func (ps *ProfileServicesImp) DeserializeUser(id string) (*models.UserProfile, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.UserProfile
	query := bson.M{"_id": oid}
	err := ps.collection.FindOne(ps.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserProfile{}, err
		}
		return nil, err
	}
	return user, nil
}
