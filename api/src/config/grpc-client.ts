import * as grpc from "@grpc/grpc-js";
import * as protoLoader from "@grpc/proto-loader";
import path from "path";

const protoPath = path.resolve(__dirname, "../../../proto/task.proto");

const packageDefinition = protoLoader.loadSync(protoPath);
const taskProto = grpc.loadPackageDefinition(packageDefinition) as any;

export const taskService = new taskProto.proto.TaskService(
  "localhost:50051",
  grpc.credentials.createInsecure()
);
