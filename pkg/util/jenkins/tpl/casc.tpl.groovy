String[] configContent = ['''[[ .CascConfig ]]''']
def configSb = new StringBuffer()
for (int i=0; i<configContent.size(); i++) {
    configSb << configContent[i]
}
def stream = new ByteArrayInputStream(configSb.toString().getBytes('UTF-8'))
def source = io.jenkins.plugins.casc.yaml.YamlSource.of(stream)
io.jenkins.plugins.casc.ConfigurationAsCode.get().configureWith(source)
