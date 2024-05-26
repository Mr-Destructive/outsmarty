export function isLoggedIn() {
  return !!document.cookie.split('; ').find(row => row.startsWith('outsmarty_uid='));
}

export function logout() {
  document.cookie = 'outsmarty_uid=; Max-Age=-99999999;';
}
