export interface GetTasksResponse {
    nextCursor: string;
    tasks: Task[];
}

export interface Task {
    id: number;
    title: string;
    content: string;
    status: 'complete' | 'in-progress' | 'incomplete';
    showActions: boolean;
}
  