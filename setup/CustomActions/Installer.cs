using System.Collections;
using System.ComponentModel;
using System.Windows.Forms;
using System.IO;
using System.Text.RegularExpressions;

namespace CustomActions
{
    [RunInstaller(true)]
    public partial class Installer : System.Configuration.Install.Installer
    {
        public Installer()
        {
            InitializeComponent();
        }

        [System.Security.Permissions.SecurityPermission(System.Security.Permissions.SecurityAction.Demand)]
        public override void Install(IDictionary stateSaver)
        {
            base.Install(stateSaver);
        }

        [System.Security.Permissions.SecurityPermission(System.Security.Permissions.SecurityAction.Demand)]
        public override void Commit(IDictionary savedState)
        {
            base.Commit(savedState);

            string type = this.Context.Parameters["type"];
            string executable = "";
            if (type.Equals("1"))
                executable = "browserpass-windows64.exe";
            if (type.Equals("2"))
                executable = "browserpass-wsl.bat";

            string path = this.Context.Parameters["targetdir"];
            path = path.Remove(path.Length - 1);
            Directory.CreateDirectory(path + "\\browser-files");
            foreach (var file in System.IO.Directory.GetFiles(path + "\\templates"))
            {
                try
                { 
                    using (StreamReader sr = new StreamReader(file))
                    {
                        string content = sr.ReadToEnd();
                        string modifiedContent = Regex.Replace(content, "path\": \"[^\"]*", "path\": \"" + path.Replace("\\", "\\\\") + executable);
                        File.WriteAllText(file.Replace("templates", "browser-files"), modifiedContent);
                    }
                }
                catch (IOException e)
                {
                    MessageBox.Show(e.Message);
                }
            }
            
        }

        [System.Security.Permissions.SecurityPermission(System.Security.Permissions.SecurityAction.Demand)]
        public override void Rollback(IDictionary savedState)
        {
            base.Rollback(savedState);
        }

        [System.Security.Permissions.SecurityPermission(System.Security.Permissions.SecurityAction.Demand)]
        public override void Uninstall(IDictionary savedState)
        {
            base.Uninstall(savedState);
            string path = this.Context.Parameters["targetdir"];
            path = path.Remove(path.Length - 1);
            Directory.Delete(path + "\\browser-files", true);
        }
    }
}
