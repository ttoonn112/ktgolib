package file

import (
  "time"
	"os"
	"strings"
	"net/http"
  "github.com/gin-gonic/gin"
)

/*
func Upload_File(c *gin.Context){

  m_user_token := c.PostForm("m_user_token")
  authorization := c.PostForm("authorization")
  folder_name := c.PostForm("folder_name")
  uid := c.PostForm("uid")
  if uid == "randomtime" {
    uid = time.Now().Format("20060102150405")+"_"+common.GenerateRandomString(5)
  }

  app_module_url := c.Request.RequestURI[1:][:strings.IndexByte(c.Request.RequestURI[1:], '/')]
  userToken := map[string]interface{}{}
  if authorization != "" {
    if user, ok := User_CheckSession(app_module_url, authorization); ok {
      userToken = User_ParseToken(user, app_module_url)
    }else{
      Msg_WithCustomlog(c, "ApiOperationError", lang.T("error.user_session_expired"), authorization)
      return
    }
  }else if m_user_token != "" {
    if user, ok := UserM_CheckSession(app_module_url, m_user_token); ok {
      userToken = user
      userToken["app_module"] = app_module_url
    }else{
      Msg_WithCustomlog(c, "ApiOperationError", lang.T("error.user_session_expired"), m_user_token)
      return
    }
  }

  if userToken == nil || common.T(userToken, "app_module") == "" {
    Msg_WithCustomlog(c, "ApiOperationError", lang.T("error.user_session_expired"), m_user_token)
    return
  }

  if _, err := os.Stat(file.File_Directory); os.IsNotExist(err) {
      os.Mkdir(file.File_Directory, os.ModePerm)
  }

  if _, err := os.Stat(file.File_Directory+"/"+userToken["app_module"].(string)); os.IsNotExist(err) {
      os.Mkdir(file.File_Directory+"/"+userToken["app_module"].(string), os.ModePerm)
  }

  filepath := file.File_Directory+"/"+userToken["app_module"].(string)+"/"+folder_name
  if _, err := os.Stat(filepath); os.IsNotExist(err) {
      os.Mkdir(filepath, os.ModePerm)
  }

  form, err := c.MultipartForm()		// Multipart form
  if err != nil {
    Msg_WithCustomlog(c, "ApiOperationError", "Upload error "+folder_name+" ("+uid+")", filepath+" : "+err.Error())
    return
  }

  fileurls := []map[string]string{}
  files := form.File["files"]
  for _, f := range files {

    filepathreal := filepath+"/"+uid+"-"+strings.Replace(f.Filename,"-","_",-1)
    common.WriteCustomLog("Upload : "+filepathreal, "FileUpload")
    if err := c.SaveUploadedFile(f, filepathreal); err != nil {
      Msg_WithCustomlog(c, "ApiOperationError", "Saving upload file error "+folder_name+" ("+uid+")", filepath+" : "+err.Error())
      return
    }

    fi, err := os.Stat(filepathreal)
    if err == nil {
      modifiedtime := strings.Replace(strings.Replace(common.DateTimeString(fi.ModTime()),":","",-1)," ","",-1)
      filepathurl :=  file.File_getUrl(userToken, folder_name+"/"+uid+"-"+strings.Replace(f.Filename,"-","_",-1)+"?"+modifiedtime)
      filepath :=  file.File_getPath(userToken, folder_name+"/"+uid+"-"+strings.Replace(f.Filename,"-","_",-1))
      thefile := map[string]string{
        "filepathurl" : filepathurl,
        "filepath" : filepath,
      }
      fileurls = append(fileurls, thefile)
    }

  }

  common.WriteCustomLog("[Upload]["+userToken["username"].(string)+"] Called", userToken["app_module"].(string)+"OperationCalled")
  c.JSON(http.StatusOK, gin.H{"what": "ok", "fileurls" : fileurls})
}*/
