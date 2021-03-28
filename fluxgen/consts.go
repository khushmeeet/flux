package fluxgen

import "path"

const TemplatesFolder = "templates"
const StaticFolder = "static"
const PagesFolder = "pages"
const SiteFolder = "_sites"
const ConfigFile = "config.yaml"

var PartialsFolder = path.Join(TemplatesFolder, "partials")
