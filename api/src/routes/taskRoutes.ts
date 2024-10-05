import { Router } from "express";
import { submitTask, getTask } from "../controllers/taskController";

const taskRouter = Router();

taskRouter.post("/", submitTask);
taskRouter.get("/", getTask);

export default taskRouter;
