// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/lebleuciel/maani/pkg/database/ent/file"
	"github.com/lebleuciel/maani/pkg/database/ent/filetype"
	"github.com/lebleuciel/maani/pkg/database/ent/schema"
	"github.com/lebleuciel/maani/pkg/database/ent/tag"
	"github.com/lebleuciel/maani/pkg/database/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	fileFields := schema.File{}.Fields()
	_ = fileFields
	// fileDescName is the schema descriptor for name field.
	fileDescName := fileFields[0].Descriptor()
	// file.NameValidator is a validator for the "name" field. It is called by the builders before save.
	file.NameValidator = func() func(string) error {
		validators := fileDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// fileDescUUID is the schema descriptor for uuid field.
	fileDescUUID := fileFields[2].Descriptor()
	// file.UUIDValidator is a validator for the "uuid" field. It is called by the builders before save.
	file.UUIDValidator = func() func(string) error {
		validators := fileDescUUID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(uuid string) error {
			for _, fn := range fns {
				if err := fn(uuid); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// fileDescType is the schema descriptor for type field.
	fileDescType := fileFields[4].Descriptor()
	// file.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	file.TypeValidator = func() func(string) error {
		validators := fileDescType.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(filetype string) error {
			for _, fn := range fns {
				if err := fn(filetype); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// fileDescCreatedAt is the schema descriptor for created_at field.
	fileDescCreatedAt := fileFields[5].Descriptor()
	// file.DefaultCreatedAt holds the default value on creation for the created_at field.
	file.DefaultCreatedAt = fileDescCreatedAt.Default.(func() time.Time)
	// fileDescUpdatedAt is the schema descriptor for updated_at field.
	fileDescUpdatedAt := fileFields[6].Descriptor()
	// file.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	file.DefaultUpdatedAt = fileDescUpdatedAt.Default.(func() time.Time)
	// file.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	file.UpdateDefaultUpdatedAt = fileDescUpdatedAt.UpdateDefault.(func() time.Time)
	// fileDescDeletedAt is the schema descriptor for deleted_at field.
	fileDescDeletedAt := fileFields[7].Descriptor()
	// file.DefaultDeletedAt holds the default value on creation for the deleted_at field.
	file.DefaultDeletedAt = fileDescDeletedAt.Default.(func() time.Time)
	filetypeFields := schema.Filetype{}.Fields()
	_ = filetypeFields
	// filetypeDescAllowedSize is the schema descriptor for allowed_size field.
	filetypeDescAllowedSize := filetypeFields[1].Descriptor()
	// filetype.DefaultAllowedSize holds the default value on creation for the allowed_size field.
	filetype.DefaultAllowedSize = filetypeDescAllowedSize.Default.(int)
	// filetypeDescIsBanned is the schema descriptor for is_banned field.
	filetypeDescIsBanned := filetypeFields[2].Descriptor()
	// filetype.DefaultIsBanned holds the default value on creation for the is_banned field.
	filetype.DefaultIsBanned = filetypeDescIsBanned.Default.(bool)
	// filetypeDescCreatedAt is the schema descriptor for created_at field.
	filetypeDescCreatedAt := filetypeFields[3].Descriptor()
	// filetype.DefaultCreatedAt holds the default value on creation for the created_at field.
	filetype.DefaultCreatedAt = filetypeDescCreatedAt.Default.(func() time.Time)
	// filetypeDescUpdatedAt is the schema descriptor for updated_at field.
	filetypeDescUpdatedAt := filetypeFields[4].Descriptor()
	// filetype.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	filetype.DefaultUpdatedAt = filetypeDescUpdatedAt.Default.(func() time.Time)
	// filetype.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	filetype.UpdateDefaultUpdatedAt = filetypeDescUpdatedAt.UpdateDefault.(func() time.Time)
	// filetypeDescID is the schema descriptor for id field.
	filetypeDescID := filetypeFields[0].Descriptor()
	// filetype.IDValidator is a validator for the "id" field. It is called by the builders before save.
	filetype.IDValidator = func() func(string) error {
		validators := filetypeDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	tagFields := schema.Tag{}.Fields()
	_ = tagFields
	// tagDescCreatedAt is the schema descriptor for created_at field.
	tagDescCreatedAt := tagFields[1].Descriptor()
	// tag.DefaultCreatedAt holds the default value on creation for the created_at field.
	tag.DefaultCreatedAt = tagDescCreatedAt.Default.(func() time.Time)
	// tagDescUpdatedAt is the schema descriptor for updated_at field.
	tagDescUpdatedAt := tagFields[2].Descriptor()
	// tag.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	tag.DefaultUpdatedAt = tagDescUpdatedAt.Default.(func() time.Time)
	// tag.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	tag.UpdateDefaultUpdatedAt = tagDescUpdatedAt.UpdateDefault.(func() time.Time)
	// tagDescID is the schema descriptor for id field.
	tagDescID := tagFields[0].Descriptor()
	// tag.IDValidator is a validator for the "id" field. It is called by the builders before save.
	tag.IDValidator = func() func(string) error {
		validators := tagDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescFirstName is the schema descriptor for first_name field.
	userDescFirstName := userFields[0].Descriptor()
	// user.FirstNameValidator is a validator for the "first_name" field. It is called by the builders before save.
	user.FirstNameValidator = func() func(string) error {
		validators := userDescFirstName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(first_name string) error {
			for _, fn := range fns {
				if err := fn(first_name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescLastName is the schema descriptor for last_name field.
	userDescLastName := userFields[1].Descriptor()
	// user.LastNameValidator is a validator for the "last_name" field. It is called by the builders before save.
	user.LastNameValidator = func() func(string) error {
		validators := userDescLastName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(last_name string) error {
			for _, fn := range fns {
				if err := fn(last_name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[2].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = userDescEmail.Validators[0].(func(string) error)
	// userDescPassword is the schema descriptor for password field.
	userDescPassword := userFields[3].Descriptor()
	// user.PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	user.PasswordValidator = userDescPassword.Validators[0].(func(string) error)
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userFields[5].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userFields[6].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
}