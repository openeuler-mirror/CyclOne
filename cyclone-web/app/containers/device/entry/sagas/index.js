let created = false;
function* empty() {
  console.log();
}

function* defaultSaga() {
  console.log();
}

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [defaultSaga];
}
