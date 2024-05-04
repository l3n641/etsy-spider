function addSpiderButton() {
    const api = 'http://127.0.0.1:8080/product'
    const div = document.querySelector('div[data-appears-component-name="search_pagination"]');
    if (div) {
        const button = document.createElement('button');
        button.textContent = '采集选中产品';
        button.style.color = 'red';
        button.style.margin = '10px';
        button.addEventListener('click', function () {
            const hrefs = getSelectLink()
            let successQuantity = 0
            for (const url of hrefs) {
                fetchSourceCode(url).then(res => {
                    console.log("请求成功链接", url)
                    successQuantity = successQuantity + 1
                    document.querySelector("#totalQuantity").textContent = "successQuantity:" + successQuantity
                    const {srcList, title} = parseHtml(res)

                    saveProduct(api, url, res, srcList, title).then(res => {
                        console.log("保存成功链接", url)
                    })
                })
            }
        });
        div.appendChild(button);

        const totalQuantity = document.createElement('span');
        totalQuantity.textContent = '选中数量:0';
        totalQuantity.setAttribute('id', 'totalQuantity');
        div.appendChild(totalQuantity);

        const successQuantity = document.createElement('span');
        successQuantity.textContent = '已经采集数:0';
        successQuantity.setAttribute('id', 'successQuantity');
        div.appendChild(successQuantity);


    }
}

// content_script.js

// Function to add checkbox to each 'link' class anchor element
function addCheckboxToLinks() {
    const links = document.querySelectorAll('ol.tab-reorder-container>li a.listing-link');

    const handleCheckboxChange = function (event) {
        const hrefs = getSelectLink()
        document.querySelector("#totalQuantity").textContent = "选中数量:" + hrefs.length
    }

    links.forEach(link => {
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.classList.add('spider-item'); // Add 'spider-item' class to checkbox
        checkbox.addEventListener('change', handleCheckboxChange)
        // Insert checkbox before the anchor element
        link.parentNode.insertBefore(checkbox, link);
    });
}

function getSelectLink() {
    const checkboxes = document.querySelectorAll('input[type="checkbox"].spider-item');

    const hrefs = [];
    checkboxes.forEach(checkbox => {
        if (checkbox.checked) {
            const href = checkbox.nextElementSibling.href;
            hrefs.push(href);
        }
    });
    return hrefs
}


//接收inject页面的消息
window.addEventListener("message", function (e) {
    if (e.data && e.data.cmd) {

    }


}, false);

//监听page_action的消息
chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {

});

window.addEventListener("load", function (event) {
    addCheckboxToLinks()
    addSpiderButton()
});


async function fetchSourceCode(url) {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            return false; // 返回 false 表示请求失败
        }
        const sourceCode = await response.text(); // 获取响应的字符串表示
        return sourceCode; // 返回源代码
    } catch (error) {
        return false; // 返回 false 表示请求失败
    }
}

async function saveProduct(api, url, htmlContent, images, title) {
    const options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({url, htmlContent, images, title})
    };

    try {
        const response = await fetch(api, options);
        return response.ok
    } catch (error) {
        return false; // 返回 false 表示请求失败
    }
}

function parseHtml(htmlText) {
    let parser = new DOMParser();
    let doc = parser.parseFromString(htmlText, 'text/html');

    let images = doc.querySelectorAll('ul[data-image-overlay-list] img');

    let srcList = [];
    images.forEach(img => {
        let src = img.getAttribute('data-src-zoom-image');
        srcList.push(src);
    });

    let title = ""
    let h1Element = doc.querySelector('[data-buy-box-listing-title="true"]');
    if (h1Element) {
        title = h1Element.textContent.trim();
    }
    return {
        srcList,
        title
    }

}



