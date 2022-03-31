package cmd

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"hitokoto-go/global"
	"hitokoto-go/models"
	"hitokoto-go/types"
	"log"
	"os"
	"path"
	"time"
)

func Import(dataDir string) {

	// Auto migrate base database
	log.Println("Start auto migrate base database...")
	if err := migrateBaseDB(); err != nil {
		log.Fatalln(err)
	}

	// Parse version
	log.Println("Start parse version...")
	versionBytes, err := os.ReadFile(path.Join(dataDir, "version.json"))
	if err != nil {
		log.Fatalln(err)
	}
	var v types.Version
	if err = json.Unmarshal(versionBytes, &v); err != nil {
		log.Fatalln(err)
	}

	// Check protocol version
	log.Println("Start check protocol version...")
	if v.ProtocolVersion != "1.0.0" {
		log.Fatalln("Unsupported protocol version: " + v.ProtocolVersion)
	}

	// Check database version
	log.Println("Start check version in database...")
	var mv models.Version
	global.DB.First(&mv)
	if mv.BundleVersion == v.BundleVersion {
		log.Println("Database is up to date")
		return
	}

	// Parse categories
	log.Println("Start parse categories...")
	categoriesBytes, err := os.ReadFile(path.Join(dataDir, v.Categories.Path))
	if err != nil {
		log.Fatalln(err)
	}
	var cs types.Categories
	if err = json.Unmarshal(categoriesBytes, &cs); err != nil {
		log.Fatalln(err)
	}

	// Save categories
	log.Println("Start save categories...")
	for _, c := range cs {
		var mc models.Category
		if err := global.DB.First(&mc, "key = ?", c.Key).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			createdAt, err := time.Parse(time.RFC3339, c.CreatedAt)
			if err != nil {
				log.Println("[WARN] Failed to parse create time for category ", c.Key, " with error: ", err)
			}
			updatedAt, err := time.Parse(time.RFC3339, c.UpdatedAt)
			if err != nil {
				log.Println("[WARN] Failed to parse update time for category ", c.Key, " with error: ", err)
			}
			mc = models.Category{
				Model: gorm.Model{
					ID:        c.ID,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				Name:        c.Name,
				Description: c.Description,
				Key:         c.Key,
			}
			global.DB.Create(&mc)
		}

		// Push to global meta, in order to create tables
		global.Meta.Categories = append(global.Meta.Categories, types.MetaCategory{
			Key: c.Key,
			// No need to initialize counts here
		})
	}

	// Create sentence tables
	log.Println("Start create sentence tables...")
	if err := migrateSentencesDB(); err != nil {
		log.Fatalln("Fail to create sentence tables with error: ", err)
	}

	// Import Sentences
	log.Println("Start import sentences...")
	for _, c := range cs {
		var ss []types.Sentence
		cssBytes, err := os.ReadFile(path.Join(dataDir, c.Path))
		if err != nil {
			log.Fatalln("Failed to load sentences of type ", c.Key, " with error: ", err)
		}
		if err = json.Unmarshal(cssBytes, &ss); err != nil {
			log.Fatalln("Failed to parse sentences of type ", c.Key, " with error: ", err)
		}
		for _, s := range ss {
			var ms models.Sentence
			if err := global.DB.
				Scopes(models.SentenceTable(models.Sentence{Type: s.Type})).
				First(&ms, "uuid = ?", s.UUID).
				Error; errors.Is(err, gorm.ErrRecordNotFound) {
				ms.FromType(&s)
				global.DB.Scopes(models.SentenceTable(ms)).Create(&ms)
			}
		}

	}

	// Save version
	log.Println("Start save version...")
	mv = models.Version{
		BundleVersion: v.BundleVersion,
	}
	global.DB.Save(&mv) // Will be auto created if not exist

	// Done!
	log.Println("All data imported! enjoy :D")

}

func migrateBaseDB() error {
	return global.DB.AutoMigrate(
		&models.Version{},
		&models.Categories{},
	)
}

func migrateSentencesDB() error {
	for _, c := range global.Meta.Categories {
		if err := global.DB.Scopes(models.SentenceTable(models.Sentence{Type: c.Key})).AutoMigrate(
			&models.Sentence{Type: c.Key},
		); err != nil {
			return err
		}
	}
	return nil
}
