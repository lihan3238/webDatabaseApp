// // 修改 searchMovies 函数
// function searchMovies() {
//     // 获取输入的用户ID
//     var userId = document.getElementById("userId").value;

//     // 创建 XMLHttpRequest 对象
//     var xhr = new XMLHttpRequest();

//     // 配置请求
//     xhr.open("POST", "/", true);
//     xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

//     // 设置回调函数
//     xhr.onreadystatechange = function () {
//         if (xhr.readyState == 4 && xhr.status == 200) {
//             // 处理返回的数据
//             displayResults(xhr.responseText);
//         }
//     };

//     // 构建请求参数
//     var params = "userId=" + userId;

//     // 发送请求
//     xhr.send(params);
// }

// // 新增函数用于显示查询结果
// function displayResults(data) {
//     // 解析 JSON 数据
//     var movies = JSON.parse(data);

//     // 获取显示结果的容器
//     var resultsContainer = document.getElementById("movieResults");

//     // 清空之前的结果
//     resultsContainer.innerHTML = "";

//     // 遍历电影数据并添加到容器中
//     for (var i = 0; i < movies.length; i++) {
//         var movie = movies[i];
//         var movieDiv = document.createElement("div");
//         movieDiv.innerHTML = "Movie: " + movie.title + ", Rating: " + movie.rating + ", Tag: " + movie.tag;
//         resultsContainer.appendChild(movieDiv);
//     }
// }
