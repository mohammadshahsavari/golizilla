package logmessages

// Define error constants
var (
	// User
	LogUserHandler    = "user_handler"
	LogUserService    = "user_service"
	LogUserRepository = "user_repository"
	LogUserModel      = "user_model"

	_                                       = ""
	LogUserLoginSuccessful                  = "user logged in successfully"
	LogUserGetProfileSuccessful             = "profile fetched successfully"
	LogUserLoginBegin                       = "starting user login process"
	LogUserVerifyLoginBegin                 = "starting user login verify process"
	LogUserVerifyLoginSuccessful            = "user login verified successfully"
	LogUserGetProfileBegin                  = "starting user Get profile"
	LogUserUpdateProfileBegin               = "starting user update profile"
	LogUserUpdateProfileSuccessful          = "user profile updated successfully"
	LogUserGetNotificationListBegin         = "starting user get notification list"
	LogUserGetNotificationListSuccessful    = "user get notification list successfully"
	LogUserGetUserByIDBegin                 = "starting user get by id"
	LogUserGetUserByIDSuccessful            = "user get by id successfully"
	LogUserGenerateAndSetJWTTokenBegin      = "starting user generate and set JWT token"
	LogUserGenerateAndSetJWTTokenSuccessful = "user generated and set JWT token successfully"
	LogUserEnable2FABegin                   = "starting user Enable 2FA"
	LogUserEnable2FASuccessful              = "user Enabled 2FA successfully"
	LogUserDisable2FABegin                  = "starting user Disable 2FA"
	LogUserDisable2FASuccessful             = "starting user Disabled 2FA successfully"
	LogUserLogoutBegin                      = "starting user logout"
	LogUserLogoutSuccessful                 = "user logged out successfully"
	LogUserVerifySignupBegin                = "starting user signup"
	LogUserVerifySignupSuccessful           = "user signed up successfully"
	LogUserCreateBegin                      = "starting create user"
	LogUserCreateSuccessful                 = "user created successfully"
	LogUserNotFound                         = "user not found"

	// Questionnaire
	LogQuestionnaireHandler                = "questionnaire_handler"
	LogQuestionnaireService                = "questionnaire_service"
	LogQuestionnaireRepository             = "questionnaire_repository"
	LogQuestionnaireModel                  = "questionnaire_model"
	_                                      = ""
	LogQuestionnaireCreateBegin            = "starting questionnaire Create"
	LogQuestionnaireDeleteBegin            = "starting questionnaire Delete"
	LogQuestionnaireUpdateBegin            = "starting questionnaire Update"
	LogQuestionnaireGetByIdBegin           = "starting questionnaire GetById"
	LogQuestionnaireGetByOwnerIdBegin      = "starting questionnaire GetByOwnerId"
	LogQuestionnaireGetResultsBegin        = "starting questionnaire GetResults"
	LogQuestionnaireGiveAccessBegin        = "starting questionnaire GiveAccess"
	LogQuestionnaireCreateSuccessful       = "questionnaire Created successfully"
	LogQuestionnaireDeleteSuccessful       = "questionnaire Deleted successfully"
	LogQuestionnaireUpdateSuccessful       = "questionnaire Updated successfully"
	LogQuestionnaireGetByIdSuccessful      = "questionnaire Got ById successfully"
	LogQuestionnaireGetByOwnerIdSuccessful = "questionnaire Got ByOwnerId successfully"
	LogQuestionnaireGiveAccessSuccessful   = "questionnaire GiveAccess successfully"
	LogQuestionnaireGetResultsEnd          = "questionnaire GetResults ended"

	// Question
	LogQuestionHandler             = "question_handler"
	LogQuestionService             = "question_service"
	LogQuestionRepository          = "question_repository"
	LogQuestionModel               = "question_model"
	_                              = ""
	LogQuestionCreateBegin         = "starting question Create"
	LogQuestionUpdateBegin         = "starting question Update"
	LogQuestionDeleteBegin         = "starting question Delete"
	LogQuestionGetByIDBegin        = "starting question GetByID"
	LogQuestionCreateSuccessfully  = "question Created successfully"
	LogQuestionUpdateSuccessfully  = "question Updated successfully"
	LogQuestionDeleteSuccessfully  = "question Deleted successfully"
	LogQuestionGetByIDSuccessfully = "question Got ByID successfully"

	// answer
	LogAnswerHandler             = "answer_handler"
	LogAnswerService             = "answer_service"
	LogAnswerRepository          = "answer_repository"
	LogAnswerModel               = "answer_model"
	_                            = ""
	LogAnswerCreateBegin         = "starting answer Create"
	LogAnswerUpdateBegin         = "starting answer Update"
	LogAnswerDeleteBegin         = "starting answer Delete"
	LogAnswerGetByIDBegin        = "starting answer GetByID"
	LogAnswerCreateSuccessfully  = "answer Created successfully"
	LogAnswerUpdateSuccessfully  = "answer Updated successfully"
	LogAnswerDeleteSuccessfully  = "answer Deleted successfully"
	LogAnswerGetByIDSuccessfully = "answer Got ByID successfully"

	// admin
	LogAdminHandler          = "Admin_handler"
	LogAdminService          = "admin_service"
	LogAdminRepository       = "admin_repository"
	LogAdminModel            = "admin_model"
	_                        = ""
	LogAdminGetAllUsersBegin = "starting admin GetAllUsers"

	LogCastUserIdError     = "failed to cast user id"
	LogLackOfAuthorization = "Lack Of Authorization"

	// role
	LogRoleService  = "roler_service"
	_               = ""
	LogRoleNotFound = "role not found"

	// role privilege on instance
	LogRolePrivilegeOnInstance = "rolePrivilegeOnInstace_repository"

	// core
	LogCoreService = "core_service"

	// submission
	LogSubmitRepo = "submission_repository"

	// Add more as needed
)
