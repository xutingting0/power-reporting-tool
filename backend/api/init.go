// file: api/init.go
package api

import "backend/repositories"

var (
	userRepo     *repositories.UserRepository
	reportRepo   *repositories.ReportRepository
	configRepo   *repositories.ConfigRepository
	templateRepo *repositories.TemplateRepository
)

// InitRepositories 初始化所有Repository
func InitRepositories(
	userRepoInstance *repositories.UserRepository,
	reportRepoInstance *repositories.ReportRepository,
	configRepoInstance *repositories.ConfigRepository,
	templateRepoInstance *repositories.TemplateRepository,
) {
	userRepo = userRepoInstance
	reportRepo = reportRepoInstance
	configRepo = configRepoInstance
	templateRepo = templateRepoInstance
}

// GetUserRepo 获取UserRepository实例
func GetUserRepo() *repositories.UserRepository {
	return userRepo
}

// GetReportRepo 获取ReportRepository实例
func GetReportRepo() *repositories.ReportRepository {
	return reportRepo
}

// GetConfigRepo 获取ConfigRepository实例
func GetConfigRepo() *repositories.ConfigRepository {
	return configRepo
}

// GetTemplateRepo 获取TemplateRepository实例
func GetTemplateRepo() *repositories.TemplateRepository {
	return templateRepo
}
