package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lebleuciel/maani/models"
	"github.com/lebleuciel/maani/pkg/database"
	"github.com/lebleuciel/maani/pkg/database/ent"
	"github.com/lebleuciel/maani/pkg/database/ent/file"
	"github.com/lebleuciel/maani/pkg/database/ent/filetype"
	"github.com/lebleuciel/maani/pkg/database/ent/migrate"
	"github.com/lebleuciel/maani/pkg/database/ent/tag"
	"github.com/lebleuciel/maani/pkg/database/ent/user"
	"github.com/pkg/errors"
)

type PostgresDatabase struct {
	db      *sql.DB
	client  *ent.Client
	baseCtx context.Context
	timeout time.Duration
}

type PostgresTransaction struct {
	PostgresDatabase
	entTx *ent.Tx
}

func (p PostgresDatabase) Migrate() error {
	err := p.client.Schema.Create(
		p.getCtx(),
		migrate.WithGlobalUniqueID(true),
	)
	if err != nil {
		return errors.Wrap(err, "Could not migrate schema to db")
	}
	return nil
}

type PGOptions struct {
	SSLMode             string
	Host                string
	User                string
	DBName              string
	Password            string
	Port                int
	MaxOpenConnections  int
	MaxIdleConnections  int
	ConnMaxLifetime     time.Duration
	Timeout             time.Duration
	ConnMaxIdleTime     time.Duration
	StatusCheckInterval time.Duration
	BaseContext         context.Context
}

func NewPostgresDatabase(options PGOptions, runMigrations bool) (*PostgresDatabase, error) {
	connectionStr := fmt.Sprintf("sslmode=%s host=%s port=%d user=%s dbname=%s password=%s",
		options.SSLMode,
		options.Host,
		options.Port,
		options.User,
		options.DBName,
		options.Password)

	db, err := sql.Open("pgx", connectionStr)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create a pgx postgres driver")
	}

	db.SetMaxOpenConns(options.MaxOpenConnections)
	db.SetMaxIdleConns(options.MaxIdleConnections)
	db.SetConnMaxLifetime(options.ConnMaxLifetime)
	db.SetConnMaxIdleTime(options.ConnMaxIdleTime)

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	pg := PostgresDatabase{
		db:      db,
		client:  client,
		baseCtx: options.BaseContext,
		timeout: options.Timeout,
	}

	if runMigrations {
		err = pg.Migrate()
		if err != nil {
			return nil, err
		}
	}

	return &pg, nil
}

func (p *PostgresDatabase) NewSerializableTransaction(ctx context.Context) (database.Transaction, error) {
	return p.NewTransaction(ctx, sql.LevelSerializable)
}

func (p *PostgresDatabase) NewTransaction(ctx context.Context, isolation sql.IsolationLevel) (database.Transaction, error) {
	entTx, err := p.client.BeginTx(ctx, &sql.TxOptions{Isolation: isolation})
	if err != nil {
		return nil, nil
	}

	var tx PostgresTransaction
	tx.PostgresDatabase = *p
	tx.client = entTx.Client()
	tx.entTx = entTx

	return &tx, nil
}

func (p *PostgresTransaction) Commit() error {
	return p.entTx.Commit()
}

func (p *PostgresTransaction) Rollback() error {
	return p.entTx.Rollback()
}

func (p *PostgresDatabase) getCtx() context.Context {
	ctx, cancel := context.WithTimeout(p.baseCtx, p.timeout)
	_ = cancel // Ignore the cancel function
	return ctx
}

func (p *PostgresDatabase) GetUserByEmail(email string) (*models.UserWithPassword, error) {
	userObj, err := p.client.User.Query().Where(user.EmailEQ(email)).Only(p.getCtx())
	if err != nil {
		return nil, errors.Wrap(err, "Could not find any users with this email address")
	}

	return &models.UserWithPassword{
		Password:    userObj.Password,
		Id:          userObj.ID,
		FirstName:   userObj.FirstName,
		LastName:    userObj.LastName,
		Email:       *userObj.Email,
		AccessType:  string(userObj.AccessType),
		CreatedAt:   *userObj.CreatedAt,
		UpdatedAt:   userObj.UpdatedAt,
		LastLoginAt: userObj.LastLoginAt,
	}, nil
}

func (p *PostgresDatabase) CreateUser(spec models.UserCreationParameters) (models.User, error) {
	_, err := p.client.User.Query().Where(user.EmailEQ(spec.Email)).Only(p.getCtx())
	var e *ent.NotFoundError
	var u *ent.User
	if errors.As(err, &e) {
		u, err = p.client.User.Create().
			SetFirstName(spec.FirstName).
			SetLastName(spec.LastName).
			SetEmail(spec.Email).
			SetPassword(spec.Password).
			SetAccessType(user.AccessType(spec.AccessType)).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			Save(p.getCtx())
		if err != nil {
			return models.User{}, errors.Wrap(err, "Could not create user")
		}
	} else {
		return models.User{}, ErrUserWithEmailExist
	}

	return models.User{
		Id:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       *u.Email,
		AccessType:  string(u.AccessType),
		CreatedAt:   *u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		LastLoginAt: u.LastLoginAt,
	}, nil
}

func (p *PostgresDatabase) UpdateUserLastLogin(userId int) error {
	_, err := p.client.User.Update().Where(user.IDEQ(userId)).SetLastLoginAt(time.Now()).Save(p.getCtx())
	return err
}

func (p *PostgresDatabase) GetUserList() ([]models.User, error) {
	users, err := p.client.User.Query().All(p.getCtx())
	if err != nil {
		return nil, err
	}

	var result []models.User
	for _, user := range users {
		result = append(result, models.User{
			Id:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       *user.Email,
			AccessType:  string(user.AccessType),
			CreatedAt:   *user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			LastLoginAt: user.LastLoginAt,
		})
	}
	return result, nil
}

func (p *PostgresDatabase) AddFileTypeIfNotExist(id string) error {
	filetype := p.client.Filetype.Query().Where(filetype.IDEQ(id)).FirstX(p.getCtx())
	if filetype == nil {
		_, err := p.client.Filetype.Create().SetID(id).SetCreatedAt(time.Now()).SetUpdatedAt(time.Now()).Save(p.getCtx())
		return err
	}
	return nil
}

func (p *PostgresDatabase) GetFileTypes() ([]models.FileType, error) {
	filetype, err := p.client.Filetype.Query().All(p.getCtx())
	if err != nil {
		return nil, err
	}

	var filetypes []models.FileType
	for _, ft := range filetype {
		filetypes = append(filetypes, models.FileType{
			Name:        ft.ID,
			AllowedSize: ft.AllowedSize,
			IsBanned:    ft.IsBanned,
		})
	}

	return filetypes, nil
}

func (p *PostgresDatabase) GetFilesSize() (int, error) {
	var sum []struct {
		Sum int
	}
	err := p.client.File.Query().Aggregate(ent.Sum(file.FieldSize)).Scan(p.getCtx(), &sum)
	if len(sum) == 0 {
		return 0, err
	}
	return sum[0].Sum, err
}

func (p *PostgresDatabase) SaveFile(file models.File) error {
	currentTagsMap := make(map[string]struct{})

	currentTags, err := p.client.Tag.Query().All(p.getCtx())
	if err != nil {
		return errors.Wrap(err, "could not get tag on saving file")
	}

	for _, t := range currentTags {
		currentTagsMap[t.ID] = struct{}{}
	}

	shouldBeAddTagsObject := make([]*ent.TagCreate, 0)
	for _, t := range file.Tags {
		if _, ok := currentTagsMap[t]; !ok {
			shouldBeAddTagsObject = append(shouldBeAddTagsObject, p.client.Tag.Create().SetID(t))
		}
	}

	err = p.client.Tag.CreateBulk(shouldBeAddTagsObject...).Exec(p.getCtx())
	if err != nil {
		return errors.Wrap(err, "could not add tag on saving file")
	}

	_, err = p.client.File.Create().
		SetName(file.Name).
		SetUUID(file.UUID).
		SetUserID(file.UserId).
		SetFiletypeID(file.TypeId).
		SetSize(file.Size).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		AddTagIDs(file.Tags...).
		Save(p.getCtx())

	return err
}

func (p *PostgresDatabase) GetFile(name []string, tags []string) (models.File, error) {
	tagQuery := p.client.Tag.Query()
	if len(tags) > 0 {
		tagQuery = tagQuery.Where(tag.IDIn(tags...))
	}
	fileQuery := tagQuery.QueryFiles()
	if len(name) > 0 {
		fileQuery = fileQuery.Where(file.NameIn(name...))
	}

	f := fileQuery.FirstX(p.getCtx())
	if f == nil {
		f = p.client.File.Query().Order(file.ByCreatedAt()).FirstX(p.getCtx())
	}
	if f == nil {
		return models.File{}, fmt.Errorf("file not found")
	}

	err := p.client.File.DeleteOneID(f.ID).Exec(p.getCtx())
	if err != nil {
		return models.File{}, err
	}

	return models.File{
		Name:   f.Name,
		UUID:   f.UUID,
		Size:   f.Size,
		TypeId: f.Type,
		UserId: f.UserID,
	}, nil
}

func (p *PostgresDatabase) GetFileList() ([]models.File, error) {
	files, err := p.client.File.Query().All(p.getCtx())
	if err != nil {
		return nil, err
	}

	var result []models.File
	for _, file := range files {
		result = append(result, models.File{
			Name:   file.Name,
			UUID:   file.UUID,
			Size:   file.Size,
			TypeId: file.Type,
			UserId: file.UserID,
		})
	}
	return result, nil
}
