'use strict';

const root = document.getElementById("root");
const _data = JSON.parse(data);

function fetchURL(url) {
  fetch(`..${url}`, {mode: 'no-cors', method: 'get'})
  .then(response => { return response.text()})
  .then(data => {
    var parser = new DOMParser();
    var doc = parser.parseFromString(data, "text/html");
    var doc_root = doc.getElementById("root");
    root.innerHTML = doc_root.innerHTML;
  })
  .catch(err => console.error(err))

  window.history.pushState('', '', url);
}

const navigation = document.querySelector('.navigation');
const links = navigation.querySelectorAll("a");

links.forEach((link) => {
  link.onclick = (e) => {
    e.preventDefault();
    if(window.location.pathname == e.target.getAttribute("href")) return
    fetchURL(e.target.getAttribute("href"));
  };
})