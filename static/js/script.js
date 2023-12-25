function searchMovies() {
    // 获取表单数据和执行查询逻辑
    const userId = document.getElementById('userId').value;
    const keyword = document.getElementById('keyword').value;
    const year = document.getElementById('year').value;
    const genre = document.getElementById('genre').value;

    // 在这里添加你的查询逻辑和与后端交互的代码
    // 使用fetch API或其他方式发送请求到后端

    // 模拟显示查询结果，实际中需要替换为从后端获取的数据
    const resultsContainer = document.getElementById('results');
    resultsContainer.innerHTML = `<p>显示查询结果的地方</p>`;
}
