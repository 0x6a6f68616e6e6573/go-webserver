'use strict';

const root = document.getElementById("root");
const _data = JSON.parse(data);

function fetchURL(url, e) {
  fetch(`../${url}`, {mode: 'no-cors', method: 'get'})
  .then(response => { return response.text()})
  .then(data => {
    var parser = new DOMParser();
    var doc = parser.parseFromString(data, "text/html");
    var doc_root = doc.getElementById("root");
    root.innerHTML = doc_root.innerHTML;
  })
  .catch(err => console.error(err))
  var target = e.target;
  if(e.target.tagName == "IMG"){
    target = e.target.parentElement;
  }
  document.querySelector('.active').classList.remove('active');
  target.classList.add("active");
  window.history.pushState('', '', url);
}

const navigation = document.querySelector('.container-fluid');
const links = navigation.querySelectorAll("a");

links.forEach((link) => {
  link.onclick = (e) => {
    e.preventDefault();
    const path = e.target.getAttribute("href") || "/";
    if(window.location.pathname == path) return;
    fetchURL(path, e);
  };
})