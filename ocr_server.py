from flask import Flask, render_template, request
import ddddocr
import requests

app = Flask(__name__)


@app.route('/ocr', methods=['GET'])
def ocr_api():
    url = request.args.get('url')
    print(url)
    html = requests.get(url,verify=False)
    ocr = ddddocr.DdddOcr()
    r = ocr.classification(html.content)
    return ''.join(r)


if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8868)
