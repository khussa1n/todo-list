package mongorepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/khussa1n/todo-list/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func (m *MongoDB) CreateTask(ctx context.Context, t *entity.Tasks) (string, error) {
	existingTaskFilter := bson.M{
		"title": t.Title,
	}

	existingTaskCount, err := m.taskCollection.CountDocuments(ctx, existingTaskFilter)
	if err != nil {
		return "", fmt.Errorf("failed to check task uniqueness: %v", err)
	}

	if existingTaskCount > 0 {
		return "", errors.New("a task with the same title and activeAt already exists")
	}

	result, err := m.taskCollection.InsertOne(ctx, t)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok == false {
		return "", errors.New("failed to convert objectid to hex")
	}

	log.Printf("create task")

	return oid.Hex(), nil
}

func (m *MongoDB) UpdateTask(ctx context.Context, t *entity.Tasks) error {
	existingTaskFilter := bson.M{
		"title": t.Title,
	}

	existingTaskCount, err := m.taskCollection.CountDocuments(ctx, existingTaskFilter)
	if err != nil {
		return fmt.Errorf("failed to check task uniqueness: %v", err)
	}

	if existingTaskCount > 0 {
		return errors.New("a task with the same title and activeAt already exists")
	}

	objectID, err := primitive.ObjectIDFromHex(t.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	taskBytes, err := bson.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal task. error: %v", err)
	}

	var updateTaskObj bson.M
	err = bson.Unmarshal(taskBytes, &updateTaskObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal task bytes. error: %v", err)
	}

	delete(updateTaskObj, "_id")

	update := bson.M{
		"$set": updateTaskObj,
	}

	result, err := m.taskCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update task. error: %v", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}

	log.Printf("update task")

	return nil
}

func (m *MongoDB) UpdateTaskStatus(ctx context.Context, id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err = m.taskCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("failed to update task status: %v", err)
	}

	return nil
}

func (m *MongoDB) GetAllTasks(ctx context.Context, status string) ([]entity.Tasks, error) {
	var tasks []entity.Tasks

	filter := bson.M{"status": status}

	findOptions := options.Find()

	cursor, err := m.taskCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks. error: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task entity.Tasks
		if err = cursor.Decode(&task); err != nil {
			return nil, fmt.Errorf("failed to decode task. error: %v", err)
		}

		tasks = append(tasks, task)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error. error: %v", err)
	}

	log.Printf("get all tasks")

	return tasks, nil
}

func (m *MongoDB) GetTaskByID(ctx context.Context, id string) (*entity.Tasks, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var task entity.Tasks

	err = m.taskCollection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, fmt.Errorf("failed to get task by ID: %v", err)
	}

	log.Printf("get task")

	return &task, nil
}

func (m *MongoDB) DeleteTask(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	result, err := m.taskCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete task. error: %v", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}

	log.Printf("delete task")

	return nil
}
