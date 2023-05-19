package viewmodels

import "go.mongodb.org/mongo-driver/bson/primitive"

type DeviceVersion struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	VersionName  string             `json:"version_name,omitempty" bson:"version_name,omitempty"`
	DeviceUuid   string             `json:"device_uuid,omitempty" bson:"device_uuid,omitempty"`
	TriedCount   int32              `json:"tried_count,omitempty" bson:"tried_count,omitempty"`
	LastModified int64              `json:"last_modified,omitempty" bson:"last_modified,omitempty"`
}

type AppVersion struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	VersionName  string             `json:"version_name,omitempty" bson:"version_name,omitempty"`
	VersionCode  string             `json:"version_code,omitempty" bson:"version_code,omitempty"`
	ApkRootPath  string             `json:"apk_root_path,omitempty" bson:"apk_root_path,omitempty"`
	ReleaseNotes []string           `json:"release_notes,omitempty" bson:"release_notes,omitempty"`
	LastModified int64              `json:"last_modified,omitempty" bson:"last_modified,omitempty"`
}

type DeviceUpdateMessage struct {
	Id                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	LatestVersion     string             `json:"latestVersion,omitempty" bson:"latestVersion,omitempty"`
	LatestVersionCode string             `json:"latestVersionCode,omitempty" bson:"latestVersionCode,omitempty"`
	Url               string             `json:"url,omitempty" bson:"url,omitempty"`
	ReleaseNotes      []string           `json:"release_notes,omitempty" bson:"release_notes,omitempty"`
}
