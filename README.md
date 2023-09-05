Hello!

I thought I would have some spare time since last friday, until now - as you may guess, I was totally wrong and i have only little time in the morning and late night :/

I wanted to be pragmatic with this task.

Making it run on colima was annoying, but I managed to find the issue :D

What could be done:

- [ ] unit/integration tests - prepared some interfaces and injected dependencies, but I am running low on personal time
- [x] logging as a middleware
- [x] healthcheck endpoint, unnecessary, but low cost effort and it is good to have it preconfigured anyway :D
- [ ] limiting the size of uploads?
- [ ] optional rate limiting, but for the size and scope of this task, I think it is not needed
- [ ] tracing - same as above, logging should be sufficient

What is unnecessary, but is here?
 - main/service.go boilerplate. I used it as a base, for now it is overkill, 
as there is not much to dispose or gracefully shutdown, yet it was a low cost effort and is not a big cognitive load to understand 
 - organisation of code might be a little over-engineered, but I wanted to play around a bit with ideas.
 - Hashing the objectID in order to create simple load balancing for the 3 instances of Minio. I think it is what you expect, based on the task, but It comes with the problems like: what if the minio node hashed for given objectID is not accessible at the moment? We return error. We can think of having duplicated data on other nodes as well, but the need of "simple load balancing" is out of the window. Similar case with having the problem being covered by replication on the Minio side. Maybe I will have some time to write test for that, at the moment only manually tested
 - http layer is a little bit leaked into the storage, but due to the time constraints I have for the next week/2, I gave up on refactoring that part for now