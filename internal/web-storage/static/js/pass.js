function confirmPasswordOnKeyUp() {
  let password = document.querySelector('#password-original');
  let confirm = document.querySelector('#password-confirm');
  let btn_submit = document.querySelector('#btn-signup-submit')
  
  if (confirm.value !== password.value) {
    btn_submit.setAttribute('disabled', "")
    password.style.borderColor = "red";
    password.style.outlineWidth = "3px";
    password.style.outline = "thick solid #f8cbd1";
    
    confirm.style.borderColor = "red";
    confirm.style.outlineWidth = "3px";
    confirm.style.outline = "thick solid #f8cbd1";

  } else {
    btn_submit.removeAttribute('disabled')
    password.style.borderColor = "#ccc";
    confirm.style.borderColor = "#ccc";

    confirm.style.outline = "none"
    password.style.outline = "none"
  }
}