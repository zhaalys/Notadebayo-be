package task

import (
	"strconv"

	"tasklybe/internal/dto"
	"tasklybe/internal/validation"
	"github.com/gofiber/fiber/v2"
)

func HandleGetTasks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	tasks, total, err := GetTasks(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"data":    nil,
			"success": false,
			"message": "Internal Server Error",
		})
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return c.JSON(dto.ResponseWrapper[[]Task]{
		Data:    tasks,
		Success: true,
		Message: "Tasks retrieved successfully",
		Pagination: &dto.PaginationResponse{
			Page:      page,
			Limit:     limit,
			Total:     int(total),
			TotalPage: totalPages,
		},
	})
}

func HandleGetDetailTask(c *fiber.Ctx) error {
	id := c.Params("id")
	task, err := GetDetailTask(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Internal Server Error",
		})
	}
	return c.JSON(dto.ResponseWrapper[Task]{
		Data:    task,
		Success: true,
		Message: "Task retrieved successfully",
	})
}

func HandleCreateTask(c *fiber.Ctx) error {

	// Disini kita deklarasi variable req yang menampung data request user
	var req CreateTaskRequest

	// Kemudian variable req divalidasi apakah request nya sesuai atau tidak
	if err := validation.BindAndValidate(c, &req); err != nil {
		// Jika tidak kita return sebagai error validation
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Validation error.",
			Error:   validation.FormatValidationError(err),
		})
	}

	// Memanggil service CreateTask
	task, err := CreateTask(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[Task]{
		Data:    task,
		Success: true,
		Message: "Success! task created.",
	})
}

func HandleEditTask(c *fiber.Ctx) error {

	/** 
	Ambil value id dari request parameter
	Contoh: http://localhost:3000/task/[xxx-xxxx-xxxx-xxxx] <- ini adalah request param
	**/
	id := c.Params("id")

	// Lakukan validasi request
	var req EditTaskRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Validation error.",
			Error:   validation.FormatValidationError(err),
		})
	}

	// Panggil service EditTask
	task, err := EditTask(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[Task]{
		Data:    task,
		Success: true,
		Message: "Success! task updated.",
	})
}

func HandleDeleteTask(c *fiber.Ctx) error {
	// Ambil id task dari url param
	id := c.Params("id")
	
	// Panggil fungsi DeleteTask, sambil di cek apakah ada error, jika ada return 500
	if err := DeleteTask(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[Task]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[Task]{
		Success: true,
		Message: "Success! task deleted.",
	})
}