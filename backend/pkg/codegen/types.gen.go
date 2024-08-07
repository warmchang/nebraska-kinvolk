// Package codegen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package codegen

import (
	"time"
)

const (
	GithubCookieAuthScopes = "githubCookieAuth.Scopes"
	OidcBearerAuthScopes   = "oidcBearerAuth.Scopes"
	OidcCookieAuthScopes   = "oidcCookieAuth.Scopes"
)

// Activity defines model for activity.
type Activity struct {
	AppID           string    `json:"app_id"`
	ApplicationName string    `json:"application_name"`
	ChannelName     string    `json:"channel_name"`
	Class           int       `json:"class"`
	CreatedTs       time.Time `db:"created_ts" json:"created_ts"`
	GroupID         string    `json:"group_id"`
	GroupName       string    `json:"group_name"`
	Id              string    `json:"id"`
	InstanceID      string    `json:"instance_id"`
	Severity        int       `json:"severity"`
	Version         string    `json:"version"`
}

// ActivityPage defines model for activityPage.
type ActivityPage struct {
	Activities []Activity `json:"activities"`
	Count      int        `json:"count"`
	TotalCount int        `json:"totalCount"`
}

// AppConfig defines model for appConfig.
type AppConfig struct {
	Description *string `json:"description,omitempty"`
	Name        string  `json:"name"`
	ProductId   *string `json:"product_id"`
}

// Application defines model for application.
type Application struct {
	Channels    []Channel `json:"channels"`
	CreatedTs   time.Time `json:"created_ts"`
	Description string    `json:"description"`
	Groups      []Group   `json:"groups"`
	Id          string    `json:"id"`
	Instances   struct {
		Count int `json:"count"`
	} `json:"instances"`
	Name      string `json:"name"`
	ProductId string `json:"product_id"`
}

// AppsPage defines model for appsPage.
type AppsPage struct {
	Applications []Application `json:"applications"`
	Count        int           `json:"count"`
	TotalCount   int           `json:"totalCount"`
}

// Arch defines model for arch.
type Arch = int

// Channel defines model for channel.
type Channel struct {
	ApplicationID string    `json:"application_id"`
	Arch          Arch      `json:"arch"`
	Color         string    `json:"color"`
	CreatedTs     time.Time `json:"created_ts"`
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Package       *Package  `json:"package,omitempty"`
	PackageID     string    `json:"package_id"`
}

// ChannelConfig defines model for channelConfig.
type ChannelConfig struct {
	ApplicationId string  `json:"application_id"`
	Arch          uint    `json:"arch"`
	Color         string  `json:"color"`
	Name          string  `json:"name"`
	PackageId     *string `json:"package_id,omitempty"`
}

// ChannelPage defines model for channelPage.
type ChannelPage struct {
	Channels   []Channel `json:"channels"`
	Count      int       `json:"count"`
	TotalCount int       `json:"totalCount"`
}

// Config defines model for config.
type Config struct {
	AccessManagementUrl string `json:"access_management_url"`
	AuthMode            string `json:"auth_mode"`
	HeaderStyle         string `json:"header_style"`
	LoginUrl            string `json:"login_url"`
	Logo                string `json:"logo"`
	LogoutUrl           string `json:"logout_url"`
	NebraskaVersion     string `json:"nebraska_version"`
	Title               string `json:"title"`
}

// ExtraFiles defines model for extraFiles.
type ExtraFiles = []struct {
	Hash    *string `json:"hash,omitempty"`
	Hash256 *string `json:"hash256,omitempty"`
	Id      *int    `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
	Size    *string `json:"size,omitempty"`
}

// FlatcarAction defines model for flatcarAction.
type FlatcarAction struct {
	ChromeOSVersion       string    `json:"chromeos_version"`
	CreatedTs             time.Time `json:"created_ts"`
	Deadline              string    `json:"deadline"`
	DisablePayloadBackoff bool      `json:"disable_payload_backoff"`
	Event                 string    `json:"event"`
	Id                    string    `json:"id"`
	IsDelta               bool      `json:"is_delta"`
	MetadataSignatureRsa  string    `json:"metadata_signature_rsa"`
	MetadataSize          string    `json:"metadata_size"`
	NeedsAdmin            bool      `json:"needs_admin"`
	Sha256                string    `json:"sha256"`
}

// FlatcarActionPackage defines model for flatcarActionPackage.
type FlatcarActionPackage struct {
	Id     *string `json:"id,omitempty"`
	Sha256 *string `json:"sha256,omitempty"`
}

// Group defines model for group.
type Group struct {
	ApplicationID             string    `json:"application_id"`
	Channel                   *Channel  `db:"channel" json:"channel,omitempty"`
	ChannelID                 string    `json:"channel_id"`
	CreatedTs                 time.Time `json:"created_ts"`
	Description               string    `json:"description"`
	Id                        string    `json:"id"`
	Name                      string    `json:"name"`
	PolicyMaxUpdatesPerPeriod int       `json:"policy_max_updates_per_period"`
	PolicyOfficeHours         bool      `json:"policy_office_hours"`
	PolicyPeriodInterval      string    `json:"policy_period_interval"`
	PolicySafeMode            bool      `json:"policy_safe_mode"`
	PolicyTimezone            string    `json:"policy_timezone"`
	PolicyUpdateTimeout       string    `json:"policy_update_timeout"`
	PolicyUpdatesEnabled      bool      `json:"policy_updates_enabled"`
	RolloutInProgress         bool      `json:"rollout_in_progress"`
	Track                     string    `json:"track"`
}

// GroupConfig defines model for groupConfig.
type GroupConfig struct {
	ChannelId                 *string `json:"channel_id,omitempty"`
	Description               *string `json:"description,omitempty"`
	Name                      string  `json:"name"`
	PolicyMaxUpdatesPerPeriod int     `json:"policy_max_updates_per_period"`
	PolicyOfficeHours         *bool   `json:"policy_office_hours,omitempty"`
	PolicyPeriodInterval      string  `json:"policy_period_interval"`
	PolicySafeMode            *bool   `json:"policy_safe_mode,omitempty"`
	PolicyTimezone            string  `json:"policy_timezone"`
	PolicyUpdateTimeout       string  `json:"policy_update_timeout"`
	PolicyUpdatesEnabled      *bool   `json:"policy_updates_enabled,omitempty"`
	Track                     *string `json:"track,omitempty"`
}

// GroupInstanceStats defines model for groupInstanceStats.
type GroupInstanceStats struct {
	Complete      int `json:"complete"`
	Downloaded    int `json:"downloaded"`
	Downloading   int `json:"downloading"`
	Error         int `json:"error"`
	Installed     int `json:"installed"`
	OnHold        int `json:"onHold"`
	Total         int `json:"total"`
	Undefined     int `json:"undefined"`
	UpdateGranted int `json:"update_granted"`
}

// GroupPage defines model for groupPage.
type GroupPage struct {
	Count      int     `json:"count"`
	Groups     []Group `json:"groups"`
	TotalCount int     `json:"totalCount"`
}

// GroupStatusCountTimeline defines model for groupStatusCountTimeline.
type GroupStatusCountTimeline = map[time.Time]map[int]map[string]uint64

// GroupVersionBreakdown defines model for groupVersionBreakdown.
type GroupVersionBreakdown = []VersionBreakdownEntry

// GroupVersionCountTimeline defines model for groupVersionCountTimeline.
type GroupVersionCountTimeline = map[time.Time]map[string]uint64

// Instance defines model for instance.
type Instance struct {
	Alias       *string              `json:"alias,omitempty"`
	Application *InstanceApplication `json:"application"`
	CreatedTs   time.Time            `json:"created_ts"`
	Id          string               `json:"id"`
	Ip          string               `json:"ip"`
}

// InstanceApplication defines model for instanceApplication.
type InstanceApplication struct {
	ApplicationID       string    `json:"application_id"`
	CreatedTs           time.Time `json:"created_ts"`
	GroupID             string    `json:"group_id"`
	InstanceID          string    `json:"instance_id"`
	LastCheckForUpdates time.Time `json:"last_check_for_updates"`
	LastUpdateGrantedTs time.Time `json:"last_update_granted_ts"`
	LastUpdateVersion   string    `json:"last_update_version"`
	Status              int       `json:"status"`
	UpdateInProgress    bool      `json:"update_in_progress"`
	Version             string    `json:"version"`
}

// InstanceCount defines model for instanceCount.
type InstanceCount struct {
	Count uint64 `json:"count"`
}

// InstancePage defines model for instancePage.
type InstancePage struct {
	Instances []Instance `json:"instances"`
	Total     int        `json:"total"`
}

// InstanceStatusHistories defines model for instanceStatusHistories.
type InstanceStatusHistories = []InstanceStatusHistory

// InstanceStatusHistory defines model for instanceStatusHistory.
type InstanceStatusHistory struct {
	CreatedTs time.Time `json:"created_ts"`
	ErrorCode string    `db:"error_code" json:"error_code"`
	Status    int       `json:"status"`
	Verison   string    `json:"verison"`
}

// LoginInfo defines model for loginInfo.
type LoginInfo struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// LoginToken defines model for loginToken.
type LoginToken struct {
	Token string `json:"token"`
}

// OmahaRequest defines model for omahaRequest.
type OmahaRequest = map[string]interface{}

// Package defines model for package.
type Package struct {
	ApplicationID     string         `json:"application_id"`
	Arch              Arch           `json:"arch"`
	ChannelsBlacklist []string       `json:"channels_blacklist"`
	CreatedTs         time.Time      `json:"created_ts"`
	Description       string         `json:"description"`
	ExtraFiles        *ExtraFiles    `json:"extra_files"`
	Filename          string         `json:"filename"`
	FlatcarAction     *FlatcarAction `json:"flatcar_action"`
	Hash              string         `json:"hash"`
	Id                string         `json:"id"`
	Size              string         `json:"size"`
	Type              int            `json:"type"`
	Url               string         `json:"url"`
	Version           string         `json:"version"`
}

// PackageConfig defines model for packageConfig.
type PackageConfig struct {
	ApplicationId     string                `json:"application_id"`
	Arch              int                   `json:"arch"`
	ChannelsBlacklist []string              `json:"channels_blacklist"`
	Description       string                `json:"description"`
	ExtraFiles        *ExtraFiles           `json:"extra_files"`
	Filename          string                `json:"filename"`
	FlatcarAction     *FlatcarActionPackage `json:"flatcar_action"`
	Hash              string                `json:"hash"`
	Size              string                `json:"size"`
	Type              int                   `json:"type"`
	Url               string                `json:"url"`
	Version           string                `json:"version"`
}

// PackagePage defines model for packagePage.
type PackagePage struct {
	Count      int       `json:"count"`
	Packages   []Package `json:"packages"`
	TotalCount int       `json:"totalCount"`
}

// UpdateInstanceConfig defines model for updateInstanceConfig.
type UpdateInstanceConfig struct {
	Alias string `json:"alias"`
}

// VersionBreakdownEntry defines model for versionBreakdownEntry.
type VersionBreakdownEntry struct {
	Instances  *int    `json:"instances,omitempty"`
	Percentage float64 `json:"percentage"`
	Version    string  `json:"version"`
}

// PaginateActivityParams defines parameters for PaginateActivity.
type PaginateActivityParams struct {
	AppIDorProductID *string `form:"appIDorProductID,omitempty" json:"appIDorProductID,omitempty"`
	GroupID          *string `form:"groupID,omitempty" json:"groupID,omitempty"`
	ChannelID        *string `form:"channelID,omitempty" json:"channelID,omitempty"`
	InstanceID       *string `form:"instanceID,omitempty" json:"instanceID,omitempty"`
	Version          *string `form:"version,omitempty" json:"version,omitempty"`
	Severity         *int    `form:"severity,omitempty" json:"severity,omitempty"`
	Start            string  `form:"start" json:"start"`
	End              string  `form:"end" json:"end"`
	Page             *int    `form:"page,omitempty" json:"page,omitempty"`
	Perpage          *int    `form:"perpage,omitempty" json:"perpage,omitempty"`
}

// PaginateAppsParams defines parameters for PaginateApps.
type PaginateAppsParams struct {
	Page    *int `form:"page,omitempty" json:"page,omitempty"`
	Perpage *int `form:"perpage,omitempty" json:"perpage,omitempty"`
}

// CreateAppParams defines parameters for CreateApp.
type CreateAppParams struct {
	CloneFrom *string `form:"clone_from,omitempty" json:"clone_from,omitempty"`
}

// PaginateChannelsParams defines parameters for PaginateChannels.
type PaginateChannelsParams struct {
	Page    *int `form:"page,omitempty" json:"page,omitempty"`
	Perpage *int `form:"perpage,omitempty" json:"perpage,omitempty"`
}

// PaginateGroupsParams defines parameters for PaginateGroups.
type PaginateGroupsParams struct {
	Page    *int `form:"page,omitempty" json:"page,omitempty"`
	Perpage *int `form:"perpage,omitempty" json:"perpage,omitempty"`
}

// GetGroupInstancesParams defines parameters for GetGroupInstances.
type GetGroupInstancesParams struct {
	Status       int     `form:"status" json:"status"`
	Page         *int    `form:"page,omitempty" json:"page,omitempty"`
	Perpage      *int    `form:"perpage,omitempty" json:"perpage,omitempty"`
	SortFilter   *string `form:"sortFilter,omitempty" json:"sortFilter,omitempty"`
	SortOrder    *string `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
	SearchFilter *string `form:"searchFilter,omitempty" json:"searchFilter,omitempty"`
	SearchValue  *string `form:"searchValue,omitempty" json:"searchValue,omitempty"`
	Duration     string  `form:"duration" json:"duration"`
	Version      *string `form:"version,omitempty" json:"version,omitempty"`
}

// GetInstanceStatusHistoryParams defines parameters for GetInstanceStatusHistory.
type GetInstanceStatusHistoryParams struct {
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`
}

// GetGroupInstanceStatsParams defines parameters for GetGroupInstanceStats.
type GetGroupInstanceStatsParams struct {
	Duration string `form:"duration" json:"duration"`
}

// GetGroupInstancesCountParams defines parameters for GetGroupInstancesCount.
type GetGroupInstancesCountParams struct {
	Duration string `form:"duration" json:"duration"`
}

// GetGroupStatusTimelineParams defines parameters for GetGroupStatusTimeline.
type GetGroupStatusTimelineParams struct {
	Duration string `form:"duration" json:"duration"`
}

// GetGroupVersionTimelineParams defines parameters for GetGroupVersionTimeline.
type GetGroupVersionTimelineParams struct {
	Duration string `form:"duration" json:"duration"`
}

// PaginatePackagesParams defines parameters for PaginatePackages.
type PaginatePackagesParams struct {
	Page          *int    `form:"page,omitempty" json:"page,omitempty"`
	Perpage       *int    `form:"perpage,omitempty" json:"perpage,omitempty"`
	SearchVersion *string `form:"searchVersion,omitempty" json:"searchVersion,omitempty"`
}

// LoginParams defines parameters for Login.
type LoginParams struct {
	LoginRedirectUrl string `form:"login_redirect_url" json:"login_redirect_url"`
}

// LoginWebhookParams defines parameters for LoginWebhook.
type LoginWebhookParams struct {
	XHubSignature string `json:"X-Hub-Signature"`
	XGithubEvent  string `json:"X-Github-Event"`
}

// CreateAppJSONRequestBody defines body for CreateApp for application/json ContentType.
type CreateAppJSONRequestBody = AppConfig

// UpdateAppJSONRequestBody defines body for UpdateApp for application/json ContentType.
type UpdateAppJSONRequestBody = AppConfig

// CreateChannelJSONRequestBody defines body for CreateChannel for application/json ContentType.
type CreateChannelJSONRequestBody = ChannelConfig

// UpdateChannelJSONRequestBody defines body for UpdateChannel for application/json ContentType.
type UpdateChannelJSONRequestBody = ChannelConfig

// CreateGroupJSONRequestBody defines body for CreateGroup for application/json ContentType.
type CreateGroupJSONRequestBody = GroupConfig

// UpdateGroupJSONRequestBody defines body for UpdateGroup for application/json ContentType.
type UpdateGroupJSONRequestBody = GroupConfig

// CreatePackageJSONRequestBody defines body for CreatePackage for application/json ContentType.
type CreatePackageJSONRequestBody = PackageConfig

// UpdatePackageJSONRequestBody defines body for UpdatePackage for application/json ContentType.
type UpdatePackageJSONRequestBody = PackageConfig

// UpdateInstanceJSONRequestBody defines body for UpdateInstance for application/json ContentType.
type UpdateInstanceJSONRequestBody = UpdateInstanceConfig

// LoginTokenFormdataRequestBody defines body for LoginToken for application/x-www-form-urlencoded ContentType.
type LoginTokenFormdataRequestBody = LoginInfo
