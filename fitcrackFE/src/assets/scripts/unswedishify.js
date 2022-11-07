import iconv from 'iconv-lite'

/* Boinc DB uses latin1_swedish encoding, which breaks names with special characters,
which are not supported in this encoding. Fix it by converting name to utf8. */
export default function (username) {
  return iconv.decode(iconv.encode(username, 'latin1'), 'utf-8')
}