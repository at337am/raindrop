// ===>>> 开始修改 1 <<<===
// 从服务器获取页面信息并触发显示。
async function fetchAndRenderPage() {
    // 状态显示的元素
    const promptElement = document.getElementById('prompt');

    try {
        const response = await fetch('/api/info');
        const info = await response.json();

        if (!response.ok) {
            throw new Error(info.error || '服务器好像开小差了');
        }

        if (info.isEmpty) {
            promptElement.textContent = "just air. and maybe some dreams.";
            return;
        }

        // 文件列表处理
        if (info.files && info.files.length > 0) {
            const fileContainer = document.getElementById('shareable-file-card');

            // 遍历文件列表, 为每个文件创建一个卡片
            info.files.forEach(file => {
                displayFileCard(file, fileContainer);
            });

            fileContainer.classList.remove('is-hidden'); // 显示整个文件容器

            // 更新页面标题
            const currentTitle = document.title;
            if (info.files.length === 1) {
                document.title = `${info.files[0].fileName} - ${currentTitle}`;
            } else {
                document.title = `${info.files.length} Files Available - ${currentTitle}`;
            }
        }

        // 描述信息处理
        if (info.description) {
            displayContent(info.description, 'description-message');
        }

        // 文本片段处理
        if (info.snippet) {
            displayContent(info.snippet, 'snippet-content');
        }

        // 最后隐藏提示元素
        promptElement.classList.add('is-hidden');
    } catch (error) {
        console.error('错误:', error);
        promptElement.textContent = `啊哦, 出了点问题! ${error.message}`;
        promptElement.classList.add('error');
    }
}

// 显示单个共享文件卡片
function displayFileCard(file, container) {
    const template = document.getElementById('file-item-template');

    // 克隆模板内容
    const fileItemClone = template.content.cloneNode(true);
    const anchorElement = fileItemClone.querySelector('.file-item');

    // 填充文件名和大小
    anchorElement.querySelector('.file-name').textContent = file.fileName;
    anchorElement.querySelector('.file-size').textContent = file.fileSize;

    // 设置正确的下载链接, 包含文件名作为查询参数
    anchorElement.href = `/api/download?file=${encodeURIComponent(file.fileName)}`;

    // 将克隆并填充好的元素添加到容器中
    container.appendChild(fileItemClone);
}

// 辅助函数, 给指定ID元素的文本内容, 并显示
function displayContent(content, elementId) {
    const element = document.getElementById(elementId);
    element.textContent = content;
    element.classList.remove('is-hidden');
}

// 从这里开始, 因为使用了 <script defer>, 可以确保 DOM 元素已加载完毕
fetchAndRenderPage();