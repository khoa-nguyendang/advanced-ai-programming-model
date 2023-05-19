# advanced-ai-programming-infrastructure
The project contains 15 services when running up with docker-compose.
This is a project for course: Programming for Artificial Intelligence.
Team members:

```
1. Nguyễn Đăng Khoa
MHV: 22C15010
khoa.nguyendang@outlook.com
--------------------
2. Phạm Minh Thạch: 
MHV: 22C15018 
thach.m.pham@gmail.com
---------------------
3. Nguyễn Y Hợp: 
MHV: 22C15006
nguyenyhop1999@gmail.com

```


prerequisites:
```
cmake
docker
docker-compose
```


prepare infrastructure

```
#prepare necessary folders
make prepare-all
```

build everything
```
make build
```

if has only trouble when pulling image `nvcr.io/nvidia/tensorrt:22.02-py3`
please manual pull it from docker (might be caused by slow internet)

```
docker pull nvcr.io/nvidia/tensorrt:22.02-py3
```


run everything
```
make up
```



### Demo
1. Enrollment.

1.1> after running all services by following above steps.
try to access `http://localhost:1080`
accept camera and you will see detection running after few second.
You may not see result , because your face did not exist in database yet.

1.2> try to access `http://localhost:1080/users`
you may see empty table at first time.
this table will contains all faces that registered in your system.
click button "add" in top right of screen

1.3> fill necessary data, and upload 5 images of yourself with 5 corners
left, right, top, bottom and straight face

1.4> if your face already exists, you may got some error returned, double check your images again.

1.5> if face is not exists, refresh `http://localhost:1080/users` to ensure your face registered in system.

2. Verify.

just back to `http://localhost:1080` and try the application.

### Information
1. Diagrams
```
README.enrollment-diagram.md : descibe enrollment flow
README.verify-diagram.md: descibe verify flow
README.structure.md: structure of project.
```