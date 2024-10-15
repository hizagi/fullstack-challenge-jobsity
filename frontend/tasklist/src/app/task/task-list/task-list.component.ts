import { Component, OnInit } from '@angular/core'
import { TaskService } from '../../services/task.service'
import { Task } from '../../models/task.model'
import { CommonModule, NgFor, NgIf } from '@angular/common'
import { FormsModule } from '@angular/forms'

@Component({
  selector: 'app-task-list',
  standalone: true,
  templateUrl: './task-list.component.html',
  styleUrls: ['./task-list.component.less'],
  imports: [NgIf, NgFor, CommonModule, FormsModule],
})
export class TaskListComponent implements OnInit {
  tasks: Task[] = []
  taskToDeleteId: number | null = null
  showDeleteModal: boolean = false
  showEditModal: boolean = false
  showCreateModal: boolean = false
  selectedTask: Task | null = null
  newTask: Task = {
    id: 0,
    title: '',
    content: '',
    status: 'incomplete',
    showActions: false,
  }

  constructor(private taskService: TaskService) {}

  ngOnInit(): void {
    this.getTasks()
  }

  getTasks(): void {
    this.taskService.getTasks().subscribe({
      next: (response) => {
        this.tasks = response.tasks.map((task) => ({
          ...task,
          showActions: false,
        }))
      },
      error: (err) => {
        console.error('Error fetching tasks:', err)
      },
    })
  }

  showActions(id: number): void {
    const task = this.tasks.find((t) => t.id === id)
    if (task) {
      task.showActions = true
    }
  }

  hideActions(id: number): void {
    const task = this.tasks.find((t) => t.id === id)
    if (task) {
      task.showActions = false
    }
  }

  openCreateModal(): void {
    this.newTask = {
      id: 0,
      title: '',
      content: '',
      status: 'incomplete',
      showActions: false,
    }
    this.showCreateModal = true
  }

  cancelCreate(): void {
    this.newTask = {
      id: 0,
      title: '',
      content: '',
      status: 'incomplete',
      showActions: false,
    }
    this.showCreateModal = false
  }

  createTask(): void {
    this.taskService.createTask(this.newTask).subscribe({
      next: () => {
        this.getTasks()
        this.cancelCreate()
      },
      error: (err) => {
        console.error('Error creating tasks:', err)
      },
    })
  }

  confirmDeleteTask(id: number): void {
    this.taskToDeleteId = id
    this.showDeleteModal = true
  }

  cancelDelete(): void {
    this.taskToDeleteId = null
    this.showDeleteModal = false
  }

  deleteTask(): void {
    if (this.taskToDeleteId !== null) {
      this.taskService.deleteTask(this.taskToDeleteId).subscribe({
        next: () => {
          this.getTasks()
          this.cancelDelete()
        },
        error: (err) => {
          console.error('Error fetching tasks:', err)
        },
      })
    }
  }

  openEditModal(task: Task): void {
    this.selectedTask = { ...task }
    this.showEditModal = true
  }

  cancelEdit(): void {
    this.selectedTask = null
    this.showEditModal = false
  }

  updateTask(): void {
    if (this.selectedTask) {
      this.taskService.updateTask(this.selectedTask).subscribe({
        next: () => {
          this.getTasks()
          this.cancelEdit()
        },
        error: (err) => {
          console.error('Error fetching tasks:', err)
        },
      })
    }
  }
}
