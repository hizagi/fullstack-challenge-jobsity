import { Component } from '@angular/core'
import { RouterOutlet } from '@angular/router'
import { TaskListComponent } from './task/task-list/task-list.component'

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, TaskListComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.less',
})
export class AppComponent {
  title = 'tasklist'
}
