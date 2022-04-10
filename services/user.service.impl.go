package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"userCrudApp/models"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func (u *UserServiceImpl) CreateUser(user *models.User) error {
	_, err := u.userCollection.InsertOne(u.ctx, user)
	return err
}

func (u *UserServiceImpl) GetUser(string2 *string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "name", Value: string2}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	result, err := u.userCollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer result.Close(u.ctx)
	for result.Next(u.ctx) {
		var user models.User
		err := result.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (u *UserServiceImpl) UpdateUser(user *models.User) error {
	filter := bson.D{bson.E{Key: "name", Value: user.Name}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "name", Value: user.Name}, bson.E{Key: "address", Value: user.Address}, bson.E{Key: "age", Value: user.Age}}}}
	result, err := u.userCollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return err
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(string2 *string) error {
	filter := bson.D{bson.E{Key: "name", Value: string2}}
	result, err := u.userCollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return err
	}
	return nil
}

func NewUserServiceImpl(userCollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		userCollection: userCollection,
		ctx:            ctx,
	}
}
