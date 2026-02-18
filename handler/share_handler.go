package handlers

import (
	"ClaranCloudDisk/model"
	services "ClaranCloudDisk/service"
	"ClaranCloudDisk/util"
	"ClaranCloudDisk/util/minIO"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShareHandler struct {
	shareService *services.ShareService
	minioClient  *minIO.MinIOClient
}

func NewShareHandler(shareService *services.ShareService, minIOClient *minIO.MinIOClient) *ShareHandler {
	return &ShareHandler{
		shareService: shareService,
		minioClient:  minIOClient,
	}
}

func (h *ShareHandler) CreateShare(c *gin.Context) {
	zap.L().Info("创建分享请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	//捕获数据
	userID := c.GetInt("user_id")
	var req model.CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorf("请求参数错误: %v", err)
		util.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 验证文件数量
	if len(req.FileIDs) == 0 {
		zap.S().Errorf("请选择要分享的文件")
		util.Error(c, 400, "请选择要分享的文件")
		return
	}

	//service
	ctx := c.Request.Context()
	share, err := h.shareService.CreateShare(ctx, uint(userID), &req)
	if err != nil {
		zap.S().Errorf("创建分享失败: %v", err)
		util.Error(c, 500, "创建分享失败: "+err.Error())
		return
	}

	// 生成分享链接
	shareURL := fmt.Sprintf("%s/share/%s", c.Request.Host, share.UniqueID)

	zap.L().Info("创建分享请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"share":       share,
		"share_url":   shareURL,
		"password":    req.Password != "",
		"expire_days": req.ExpireDays,
		"expire_time": share.CreatedAt.Add(share.Exp).Format("2006-01-02 15:04:05"),
	}, "分享创建成功")
}

func (h *ShareHandler) CheckMine(c *gin.Context) {
	zap.L().Info("查看个人分享请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	userID := c.GetInt("user_id")

	ctx := c.Request.Context()
	shares, total, err := h.shareService.GetMyShares(ctx, uint(userID))
	if err != nil {
		zap.S().Errorf("获取分享列表失败: %v", err)
		util.Error(c, 500, "获取分享列表失败: "+err.Error())
		return
	}

	zap.L().Info("查看个人分享请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"shares": shares,
		"total":  total,
	}, "获取成功")
}

func (h *ShareHandler) DeleteShare(c *gin.Context) {
	zap.L().Info("删除个人分享请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	userID := c.GetInt("user_id")
	uniqueID := c.Param("unique_id")

	ctx := c.Request.Context()
	err := h.shareService.DeleteShare(ctx, uint(userID), uniqueID)
	if err != nil {
		zap.S().Errorf("删除分享失败: %v", err)
		util.Error(c, 403, err.Error())
		return
	}

	zap.L().Info("删除个人分享请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{}, "删除分享成功")
}

func (h *ShareHandler) GetShareInfo(c *gin.Context) {
	zap.L().Info("获取分享信息请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	uniqueID := c.Param("unique_id")
	password := c.PostForm("password")
	//zap.S().Info(password)

	ctx := c.Request.Context()
	shareInfo, err := h.shareService.GetShareInfo(ctx, uniqueID, password)
	if err != nil {
		zap.S().Errorf("获取分享信息失败: %v", err)
		util.Error(c, 403, err.Error())
		return
	}

	// 生成完整的分享链接
	shareURL := fmt.Sprintf("%s/share/%s", c.Request.Host, uniqueID)

	zap.L().Info("获取分享信息请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"share":         shareInfo.Share,
		"files":         shareInfo.Files,
		"need_password": shareInfo.NeedPassword,
		"is_expired":    shareInfo.IsExpired,
		"expire_time":   shareInfo.ExpireTime,
		"share_url":     shareURL,
		"total_size":    shareInfo.TotalSize,
		"file_count":    shareInfo.FileCount,
	}, "获取分享信息成功")
}

func (h *ShareHandler) DownloadSpecFile(c *gin.Context) {
	zap.L().Info("下载特定文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	userID := c.GetInt("user_id")
	uniqueID := c.Param("unique_id")
	fileIDStr := c.Param("file_id")
	password := c.DefaultQuery("password", "")

	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	ctx := c.Request.Context()
	file, limitedSpeed, err := h.shareService.DownloadSpecFile(ctx, uniqueID, password, uint(fileID), userID)
	if err != nil {
		zap.S().Errorf("下载指定文件失败: %v", err)
		util.Error(c, 403, err.Error())
		return
	}

	//// 设置下载响应头
	//c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	//c.Header("Content-Type", "application/octet-stream")
	//c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
	//
	//zap.L().Info("下载特定文件请求结束",
	//	zap.String("url", c.Request.RequestURI),
	//	zap.String("method", c.Request.Method),
	//	zap.String("client_ip", c.ClientIP()))
	//
	//// 发送文件
	//c.File(file.Path)
	//设置响应头，返回的信息为下载文件流本身，而非JSON响应
	//指定传输编码为二进制，确保文件不会因为编码问题而损坏
	c.Header("Content-Transfer-Encoding", "binary")
	//强制下载并指定文件名
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	//设置文件类型为二进制文件
	c.Header("Content-Type", "application/octet-stream")
	//提供Size用于为客户端提供下载进度和剩余时间
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))

	//从minIO获取文件流
	stream, err := h.minioClient.GetStream(c, file.Path)
	if err != nil {
		zap.S().Errorf("从minIO获取文件失败: %v", err)
		util.Error(c, 500, "从minIO获取文件失败"+err.Error())
		return
	}
	defer stream.Close()

	//不限速
	if limitedSpeed == 0 {
		io.Copy(c.Writer, stream)
		zap.L().Info("下载请求结束",
			zap.String("url", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()))
		return
	}

	//限速
	bufferSize := int64(64 * 1024) // 64KB缓冲区
	if limitedSpeed < bufferSize {
		bufferSize = limitedSpeed
	}

	buf := make([]byte, bufferSize)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 每秒最多读取limitedSpeed字节
			bytesRead := int64(0)
			for bytesRead < limitedSpeed {
				remaining := limitedSpeed - bytesRead
				readSize := remaining
				if readSize > bufferSize {
					readSize = bufferSize
				}

				// 读取文件
				n, err := stream.Read(buf[:readSize])
				if n > 0 {
					// 写入HTTP响应
					_, writeErr := c.Writer.Write(buf[:n])
					if writeErr != nil {
						return
					}
					c.Writer.Flush()      // 立即发送给客户端
					bytesRead += int64(n) // 累计已读取字节
				}

				if err != nil {
					if err == io.EOF {
						zap.L().Info("下载请求结束",
							zap.String("url", c.Request.RequestURI),
							zap.String("method", c.Request.Method),
							zap.String("client_ip", c.ClientIP()))
						return // 文件读取完成
					}
					return
				}
			}
		case <-ctx.Done():
			zap.L().Info("下载请求超时",
				zap.String("url", c.Request.RequestURI),
				zap.String("method", c.Request.Method),
				zap.String("client_ip", c.ClientIP()))
			return // 上下文取消
		}
	}

	//=============================================================================================================
	//发送文件
	//fileContent, err := os.Open(file.Path)
	//if err != nil {
	//	util.Error(c, 500, "打开文件失败: "+err.Error())
	//	return
	//}
	//defer fileContent.Close()
	//
	//// 不限速
	//if limitedSpeed == 0 {
	//	io.Copy(c.Writer, fileContent)
	//	return
	//}
	//
	//// 限速处理
	//bufferSize := int64(64 * 1024) // 64KB缓冲区
	//if limitedSpeed < bufferSize {
	//	bufferSize = limitedSpeed
	//}
	//
	//buf := make([]byte, bufferSize)
	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()
	//
	//for {
	//	select {
	//	case <-ticker.C:
	//		// 每秒最多读取limitedSpeed字节
	//		bytesRead := int64(0)
	//		for bytesRead < limitedSpeed {
	//			remaining := limitedSpeed - bytesRead
	//			readSize := remaining
	//			if readSize > bufferSize {
	//				readSize = bufferSize
	//			}
	//
	//			// 读取文件
	//			n, err := fileContent.Read(buf[:readSize])
	//			if n > 0 {
	//				// 写入HTTP响应
	//				_, writeErr := c.Writer.Write(buf[:n])
	//				if writeErr != nil {
	//					return
	//				}
	//				c.Writer.Flush()      // 立即发送给客户端
	//				bytesRead += int64(n) // 累计已读取字节
	//			}
	//
	//			if err != nil {
	//				if err == io.EOF {
	//					return // 文件读取完成
	//				}
	//				return
	//			}
	//		}
	//	case <-ctx.Done():
	//		return // 上下文取消
	//	}
	//}
	//=============================================================================================================
}

func (h *ShareHandler) SaveSpecFile(c *gin.Context) {
	zap.L().Info("转存特定文件请求开始",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))
	userID := c.GetInt("user_id")
	uniqueID := c.Param("unique_id")
	fileIDStr := c.Param("file_id")
	password := c.DefaultQuery("password", "")

	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		zap.S().Errorf("无效的文件ID: %v", err)
		util.Error(c, 400, "无效的文件ID")
		return
	}

	ctx := c.Request.Context()
	savedFile, err := h.shareService.SaveSpecFile(ctx, uint(userID), uniqueID, password, uint(fileID))
	if err != nil {
		zap.S().Errorf("转存文件失败: %v", err)
		util.Error(c, 500, "转存文件失败: "+err.Error())
		return
	}

	zap.L().Info("转存特定文件请求结束",
		zap.String("url", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()))

	util.Success(c, gin.H{
		"file": savedFile,
	}, "文件转存成功")
}
