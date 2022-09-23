import jenkins.model.*
import java.util.logging.Logger
import hudson.util.VersionNumber

def logger = Logger.getLogger("")
def installed = false
def initialized = false
def pluginParameter="[[ .JenkinsPlugins ]]"
def plugins = pluginParameter.split()
def instance = Jenkins.getInstance()
def pm = instance.getPluginManager()
def uc = instance.getUpdateCenter()
// install plugins
plugins.each {
  // get plugin name and version
  String[] pluginString = it.split(":")
  String pluginName = pluginString[0]
  String pluginVersion = pluginString[1]
  logger.info("Checking " + pluginName)
  if (!pm.getPlugin(pluginName)) {
    logger.info("Looking UpdateCenter for " + pluginName)
    if (!initialized) {
      // update jenkins update center
      uc.updateAllSites()
      initialized = true
    }
    def plugin = uc.getPlugin(pluginName, new VersionNumber(pluginVersion))
    if (plugin) {
      // install plugin
      logger.info("Installing " + pluginName)
      def installFuture = plugin.deploy()
      while(!installFuture.isDone()) {
        logger.info("Waiting for plugin install: " + pluginName)
        // wait 1 second for plugin installtion
        sleep(1000)
      }
      installed = true
    } else {
      logger.warn("Plugin version not exist")
    }
  }
}
if (installed) {
  logger.info("Plugins installed, initializing a restart!")
  instance.save()

  [[ if .EnableRestart ]]
  instance.doSafeRestart()
  [[ end ]]
}
