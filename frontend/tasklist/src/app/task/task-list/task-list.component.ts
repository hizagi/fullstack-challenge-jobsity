import { Component, OnInit } from '@angular/core';
import { TaskService } from '../services/task.service';
import { Task } from '../models/task.model';
import { NgFor } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-task-list',
  standalone: true,
  templateUrl: './task-list.component.html',
  styleUrls: ['./task-list.component.less'],
  imports: [NgFor, FormsModule],
})
export class TaskListComponent implements OnInit {
  tasks: Task[] = [];
  taskToDeleteId: number | null = null;
  showDeleteModal: boolean = false;

  constructor(private taskService: TaskService) {}

  ngOnInit(): void {
    this.getTasks();
  }

  getTasks(): void {
    this.taskService.getTasks().subscribe((tasks) => {
      this.tasks = tasks.map((task) => ({
        ...task,
        showActions: false, // For handling hover display
      }));
    });
  }

  showActions(id: number): void {
    const task = this.tasks.find((t) => t.id === id);
    if (task) {
      task.showActions = true;
    }
  }

  hideActions(id: number): void {
    const task = this.tasks.find((t) => t.id === id);
    if (task) {
      task.showActions = false;
    }
  }

  confirmDeleteTask(id: number): void {
    this.taskToDeleteId = id;
    this.showDeleteModal = true;
  }

  cancelDelete(): void {
    this.taskToDeleteId = null;
    this.showDeleteModal = false;
  }

  deleteTask(): void {
    if (this.taskToDeleteId !== null) {
      this.taskService.deleteTask(this.taskToDeleteId).subscribe(() => {
        this.getTasks(); // Reload tasks after deletion
        this.cancelDelete(); // Close modal
      });
    }
  }

  editTask(task: Task): void {
    // Implement the logic for editing a task
    console.log('Edit task:', task);
  }
}
