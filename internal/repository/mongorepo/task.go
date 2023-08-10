package mongorepo

import (
	"context"
	"fmt"
	"github.com/khussa1n/todo-list/internal/custom_error"
	"github.com/khussa1n/todo-list/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func (m *MongoDB) CreateTask(ctx context.Context, t *entity.Tasks) (*entity.Tasks, error) {
	existingTaskFilter := bson.M{
		"title": t.Title,
	}

	existingTaskCount, err := m.taskCollection.CountDocuments(ctx, existingTaskFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to check task uniqueness: %v", err)
	}

	if existingTaskCount > 0 {
		return nil, custom_error.ErrDuplicateTask
	}

	result, err := m.taskCollection.InsertOne(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %v", err)
	}

	t.ID = result.InsertedID.(primitive.ObjectID)

	log.Printf("create task")

	return t, nil
}

func (m *MongoDB) UpdateTask(ctx context.Context, t *entity.Tasks, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

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
		return custom_error.ErrTaskNotFound
	}

	log.Printf("update task")

	return nil
}

func (m *MongoDB) UpdateTaskStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := m.taskCollection.UpdateOne(ctx, filter, update)
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

func (m *MongoDB) GetTaskByID(ctx context.Context, id primitive.ObjectID) (*entity.Tasks, error) {
	filter := bson.M{"_id": id}

	var task entity.Tasks

	err := m.taskCollection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, custom_error.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task by ID: %v", err)
	}

	log.Printf("get task")

	return &task, nil
}

func (m *MongoDB) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := m.taskCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete task. error: %v", err)
	}

	if result.DeletedCount == 0 {
		return custom_error.ErrTaskNotFound
	}

	log.Printf("delete task")

	return nil
}
