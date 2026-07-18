const source = await Bun.file("vpkg/vandor/go/templates/http_handler.vxt").text();

for (const fragment of [
  "type input struct{}",
  "type output struct{}",
  'method.Tags("{{ path }}")',
]) {
  if (!source.includes(fragment)) throw new Error(`HTTP handler template is missing ${fragment}`);
}
