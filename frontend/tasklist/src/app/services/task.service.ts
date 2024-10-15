import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders, HttpParams } from '@angular/common/http';
import { catchError, Observable, throwError } from 'rxjs';
import { GetTasksResponse, Task } from '../models/task.model';

@Injectable({
  providedIn: 'root',
})
export class TaskService {
  private apiUrl = `${import.meta.env.NG_APP_BASE_URL}/tasks`;

  constructor(private http: HttpClient) {}

  getTasks(cursor: string = '', limit: number = 10): Observable<GetTasksResponse> {
    const headers = new HttpHeaders({
      'Authorization': `${import.meta.env.NG_APP_API_KEY}`,  // Add the API key in the headers
    });

    let params = new HttpParams().set('limit', limit);
    params = params.set('cursor', '');
    

    return this.http.get<GetTasksResponse>(this.apiUrl, { headers, params }).pipe(
      catchError(this.handleError)  // Handle errors
    );
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

    return this.http.patch<Task>(`${this.apiUrl}/${task.id}`, task, { headers });
  }

  deleteTask(id: number): Observable<void> {
    const headers = new HttpHeaders({
      'Authorization': `${import.meta.env.NG_APP_API_KEY}`,
    });

    return this.http.delete<void>(`${this.apiUrl}/${id}`, { headers });
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    let errorMessage = '';
    if (error.error instanceof ErrorEvent) {
      // Client-side error
      errorMessage = `Client-side error: ${error.error.message}`;
    } else {
      // Server-side error
      errorMessage = `Server error (status: ${error.status}): ${error.message}`;
    }
    console.error(errorMessage);
    return throwError(() => new Error(errorMessage));  // Throw the error to be handled by the component
  }
}
