package constants

var ROLE_TPR = "TPR"
var ROLE_STUDENT = "STUDENT"
var ROLE_RECRUITER = "recruiter" // it is not capitalized in DB
var ROLE_ADMIN = "ADMIN"

var ROLE_GROUP_READ = "GROUP_READ"
var ROLE_GROUP_EDIT = "GROUP_EDIT"
var ROLE_GROUP_CREATE = "GROUP_CREATE"
var ROLE_GROUP_DELETE = "GROUP_DELETE"

var ENV_STUDENT_GROUP_OBJ_ID = "STUDENT_GROUP_OBJ_ID"

type Action string

const (
	ACTION_PUSH Action = "PUSH"
	ACTION_PULL Action = "PULL"
)
