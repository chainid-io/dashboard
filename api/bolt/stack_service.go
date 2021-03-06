package bolt

import (
	"github.com/chainid-io/dashboard"
	"github.com/chainid-io/dashboard/bolt/internal"

	"github.com/boltdb/bolt"
)

// StackService represents a service for managing stacks.
type StackService struct {
	store *Store
}

// Stack returns a stack object by ID.
func (service *StackService) Stack(ID chainid.StackID) (*chainid.Stack, error) {
	var data []byte
	err := service.store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(stackBucketName))
		value := bucket.Get([]byte(ID))
		if value == nil {
			return chainid.ErrStackNotFound
		}

		data = make([]byte, len(value))
		copy(data, value)
		return nil
	})
	if err != nil {
		return nil, err
	}

	var stack chainid.Stack
	err = internal.UnmarshalStack(data, &stack)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

// Stacks returns an array containing all the stacks.
func (service *StackService) Stacks() ([]chainid.Stack, error) {
	var stacks = make([]chainid.Stack, 0)
	err := service.store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(stackBucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var stack chainid.Stack
			err := internal.UnmarshalStack(v, &stack)
			if err != nil {
				return err
			}
			stacks = append(stacks, stack)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return stacks, nil
}

// StacksBySwarmID return an array containing all the stacks related to the specified Swarm ID.
func (service *StackService) StacksBySwarmID(id string) ([]chainid.Stack, error) {
	var stacks = make([]chainid.Stack, 0)
	err := service.store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(stackBucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var stack chainid.Stack
			err := internal.UnmarshalStack(v, &stack)
			if err != nil {
				return err
			}
			if stack.SwarmID == id {
				stacks = append(stacks, stack)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return stacks, nil
}

// CreateStack creates a new stack.
func (service *StackService) CreateStack(stack *chainid.Stack) error {
	return service.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(stackBucketName))

		data, err := internal.MarshalStack(stack)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(stack.ID), data)
		if err != nil {
			return err
		}
		return nil
	})
}

// UpdateStack updates an stack.
func (service *StackService) UpdateStack(ID chainid.StackID, stack *chainid.Stack) error {
	data, err := internal.MarshalStack(stack)
	if err != nil {
		return err
	}

	return service.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(stackBucketName))
		err = bucket.Put([]byte(ID), data)
		if err != nil {
			return err
		}
		return nil
	})
}

// DeleteStack deletes an stack.
func (service *StackService) DeleteStack(ID chainid.StackID) error {
	return service.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(stackBucketName))
		err := bucket.Delete([]byte(ID))
		if err != nil {
			return err
		}
		return nil
	})
}
