import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Task } from '../models/task.model';

@Injectable({
  providedIn: 'root',
})
export class TaskService {
  private apiUrl = `${import.meta.env.NG_APP_BASE_URL}/tasks`;

  constructor(private http: HttpClient) {}

  getTasks(cursor: string = '', limit: number = 10): Observable<Task[]> {
    const headers = new HttpHeaders({
      'Authorization': `${import.meta.env.NG_APP_API_KEY}`,  // Add the API key in the headers
    });

    let params = new HttpParams().set('limit', limit);
    if (cursor) {
      params = params.set('cursor', cursor);
    }

    return this.http.get<Task[]>(this.apiUrl, { headers, params });
  }

  getTask(id: number): Observable<Task> {
    const headers = new HttpHeaders({
      'Authorization': `${import.meta.env.NG_APP_API_KEY}`,
    });

    return this.http.get<Task>(`${this.apiUrl}/${id}`, { headers });
  }

  createTask(task: Task): Observable<Task> {
    const headers = new HttpHeaders({
      'Authorization': `${import.meta.env.NG_APP_API_KEY}`,
    });

    return this.http.post<Task>(this.apiUrl, task, { headers });
  }

  updateTask(task: Task): Observable<Task> {
    const headers = new HttpHeaders({
      'Authorization': `${import.meta.env.NG_APP_API_KEY}`,
    });

    return this.http.put<Task>(`${this.apiUrl}/${task.id}`, task, { headers });
  }

  deleteTask(id: number): Observable<void> {
    const headers = new HttpHeaders({
      'Authorization': `${import.meta.env.NG_APP_API_KEY}`,
    });

    return this.http.delete<void>(`${this.apiUrl}/${id}`, { headers });
  }
}
