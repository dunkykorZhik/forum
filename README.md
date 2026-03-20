# Forum App

**Author:** Overbuy Korkemay  

A simple forum web application where users can create posts, comment, like/dislike, and interact with categories. Admins and moderators can manage content and handle reports.  

**Built with:**  
- Go (standard packages)  
- SQLite3 (database)  
- bcrypt (password hashing)  
- github.com/google/uuid (unique IDs)  

---

## How to Run

### Using Docker

1. Build and start the container:

```bash
docker-compose up --build
```

2. The Server is on
```bash
http://localhost:8080
```

*for different docker version try 
```bash
docker compose build
docker compose up
```