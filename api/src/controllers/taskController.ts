import { Request, Response } from "express";
import { taskService } from "../config/grpc-client";
import { Task, TaskResponse } from "../types/tasks";

export const submitTask = async (req: Request, res: Response) => {
  const task: Task = req.body;
  taskService.SubmitTask(task, (err: Error | null, response: TaskResponse) => {
    if (err) {
      res.status(500).json({ error: err.message });
    } else {
      res.json(response);
    }
  });
};

export const getTask = async (req: Request, res: Response) => {
  taskService.GetTask({}, (err: Error | null, task: Task) => {
    if (err) {
      res.status(500).json({ error: err.message });
    } else {
      res.json(task);
    }
  });
};
