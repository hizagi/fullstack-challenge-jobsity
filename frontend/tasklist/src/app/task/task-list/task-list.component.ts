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
  nextCursor: string | null = null
  limit: number = 10
  isLoading: boolean = false
  hasMoreTasks: boolean = true

  constructor(private taskService: TaskService) {}

  ngOnInit(): void {
    this.loadTasks()
  }

  resetAndLoad(): void {
    this.nextCursor = ''
    this.tasks = []
    this.loadTasks()
  }

  loadTasks(): void {
    if (this.isLoading) {
      return
    }

    this.isLoading = true

    this.taskService.getTasks(this.nextCursor || '', this.limit).subscribe({
      next: (response) => {
        if (response.tasks?.length > 0) {
          this.tasks = [...this.tasks, ...response.tasks]
        }

        this.nextCursor = response.nextCursor || null
        this.hasMoreTasks = !!response.nextCursor

        this.isLoading = false
      },
      error: (err) => {
        console.error('Error fetching tasks:', err)
        this.isLoading = false
      },
    })
  }

  loadMoreTasks(): void {
    this.loadTasks()
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
        this.resetAndLoad()
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
          this.resetAndLoad()
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
          this.resetAndLoad()
          this.cancelEdit()
        },
        error: (err) => {
          console.error('Error fetching tasks:', err)
        },
      })
    }
  }

  onLimitChange(newLimit: number): void {
    this.limit = newLimit
    this.resetAndLoad()
  }
}
