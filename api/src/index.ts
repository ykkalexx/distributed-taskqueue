import express from "express";
import bodyParser from "body-parser";
import cors from "cors";
import taskRouter from "./routes/taskRoutes";

const app = express();
const port = 3000;

app.use(cors());
app.use(bodyParser.json());

app.use("/api/tasks", taskRouter);

app.listen(port, () => {
  console.log(`Server is running at http://localhost:${port}`);
});
