# High-level Flow

Scheduler calls the ReportJob
Job get pushed to jobQueue
WorkerPool recieves the job, calls the executor.Execute(job)
Executor then loads, template render, outputGenerator and then calls the Delivery adapter that also handles logs
