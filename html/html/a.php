<?php
date_default_timezone_set('Asia/Shanghai');
require 'vendor/autoload.php';

// 获取请求的日期参数
$requestedDate = isset($_GET['date']) ? $_GET['date'] : date('Ymd');
$baseDir = '/robotv1/articles/';
$filename = $baseDir . $requestedDate . '.md';

// 检查是否请求文章列表
$showList = isset($_GET['action']) && $_GET['action'] === 'list';

// 扫描文章目录
$articles = [];
if ($handle = opendir($baseDir)) {
    while (false !== ($entry = readdir($handle))) {
        if (preg_match('/^(\d{8})\.md$/', $entry, $matches)) {
            $date = $matches[1];
            $dateObj = DateTime::createFromFormat('Ymd', $date);
            $formattedDate = $dateObj ? $dateObj->format('Y年m月d日') : $date;
            $articles[$date] = [
                'filename' => $entry,
                'display_date' => $formattedDate,
                'filepath' => $baseDir . $entry
            ];
        }
    }
    closedir($handle);
    krsort($articles); // 按日期降序排列
}

echo '<!DOCTYPE html><html><head><meta charset="UTF-8">
<title>文档浏览 - ' . $requestedDate . '</title>
<style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; 
            max-width: 800px; margin: 20px auto; padding: 20px; }
    .date-selector { background: #f5f5f5; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
    .date-selector .form-container { display: flex; align-items: center; }
    .date-selector input, .date-selector .ad-button { padding: 8px 12px; font-size: 16px; }
    .date-selector input { flex: 1; margin-right: 10px; }
    .ad-button { background: #e74c3c; color: white; border: none; border-radius: 4px; 
                cursor: pointer; text-decoration: none; display: inline-block; }
    .ad-button:hover { background: #c0392b; }
    .error-container { text-align: center; padding: 40px 0; }
    .error-container h1 { color: #e74c3c; }
    .nav-links { margin: 15px 0; }
    .nav-links a { display: inline-block; margin-right: 10px; padding: 5px 10px; 
                  background: #3498db; color: white; border-radius: 4px; text-decoration: none; }
    .nav-links a:hover { background: #2980b9; }
    .article-list { margin-top: 20px; }
    .article-item { padding: 10px; border-bottom: 1px solid #eee; }
    .article-item:hover { background: #f9f9f9; }
    .article-date { font-weight: bold; color: #2c3e50; }
    .back-link { display: block; margin: 15px 0; }
    
    /* 其他原有样式保持不变 */
    h1, h2, h3 { color: #2c3e50; border-bottom: 1px solid #eee; padding-bottom: 10px; }
    code { background: #f8f8f8; padding: 2px 5px; border-radius: 3px; }
    pre { background: #f8f8f8; padding: 10px; border-radius: 5px; overflow: auto; }
    blockquote { border-left: 4px solid #ddd; padding-left: 15px; color: #777; }
    a { color: #3498db; text-decoration: none; }
    a:hover { text-decoration: underline; }
</style>
</head><body>';

// 日期选择器
echo '<div class="date-selector">
    <div class="form-container">
        <form method="GET" action="" style="flex:1; display:flex;">
            <input type="date" id="dateInput" name="date" 
                   value="' . substr($requestedDate, 0, 4) . '-' . 
                          substr($requestedDate, 4, 2) . '-' . 
                          substr($requestedDate, 6, 2) . '">
        </form>
        <!-- 广告位招租按钮放在日期选择框右侧 -->
        <a href="https://ttpai.cn" target="_blank" class="ad-button">广告位招租</a>
    </div>
    <div class="nav-links">
        <a href="?date=' . date('Ymd', strtotime($requestedDate . ' -1 day')) . '">← 前一天</a>
        <a href="?date=' . date('Ymd') . '">今天</a>
        <a href="?date=' . date('Ymd', strtotime($requestedDate . ' +1 day')) . '">后一天 →</a>
        <a href="?action=list" style="background:#27ae60;">查看历史文章</a>
    </div>
</div>';

// 显示文章列表
if ($showList) {
    echo '<h1>历史文章列表</h1>';
    echo '<div class="article-list">';
    
    if (!empty($articles)) {
        foreach ($articles as $date => $article) {
            echo '<div class="article-item">';
            echo '<span class="article-date">' . $article['display_date'] . '</span>';
            echo '<a href="?date=' . $date . '">查看文档</a>';
            echo '</div>';
        }
    } else {
        echo '<p>未找到任何文章</p>';
    }
    
    echo '</div>';
    echo '<a class="back-link" href=".">返回今日文章</a>';
} 
// 显示单篇文章
else {
    if (file_exists($filename)) {
        $markdown = file_get_contents($filename);
        $parsedown = new Parsedown();
        $html = $parsedown->text($markdown);
        echo '<h1>' . $requestedDate . ' 文档</h1>';
        echo $html;
    } else {
        $formattedDate = date_create_from_format('Ymd', $requestedDate)->format('Y年m月d日');
        echo '<div class="error-container">
            <h1>未找到文档</h1>
            <p>未找到 ' . $formattedDate . ' 的文档（' . basename($filename) . '）</p>
            <p>请确认以下事项：</p>
            <ul style="text-align: left; max-width: 400px; margin: 20px auto;">
                <li>文档是否已创建在正确目录</li>
                <li>文件名格式是否为 YYYYMMDD.md</li>
                <li>当前请求日期：' . $formattedDate . '</li>
                <li>联系电话：17678601514</li>
            </ul>
            <p><a href="?action=list">查看所有历史文章</a></p>
        </div>';
    }
}

// 高亮JS
echo '<link rel="stylesheet" href="github.min.css">
      <script src="highlight.min.js"></script>
      <script>hljs.highlightAll();</script>
      </body></html>';
?>
