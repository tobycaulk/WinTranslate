using Google.Api.Gax.ResourceNames;
using Google.Cloud.Translate.V3;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;

namespace WinTranslate
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        private static readonly string LANGUAGE_LABEL_TEMPLATE = "{0} -> {1}";
        private static readonly string TranslationServiceUrl = "http://localhost:3000/translate";

        private TranslationServiceClient googleCloudClient;
        private bool translateToSwedish = true;

        private class TranslationRequest
        {
            public string Text { get; set; }
            public string LanguageCode { get; set; }
        }

        private class TranslationResponse
        {
            public string Text { get; set; }
        }

        public MainWindow()
        {
            this.Loaded += Window_Loaded;

            InitializeComponent();

            TranslationServiceClientBuilder builder = new TranslationServiceClientBuilder
            {
                CredentialsPath = "C:\\Users\\tobyc\\Documents\\Git\\Environment\\wintranslate-api-key.json"
            };
            googleCloudClient = builder.Build();
        }

        private void Window_Loaded(object sender, RoutedEventArgs e)
        {
            var desktopWorkingArea = SystemParameters.WorkArea;
            this.Left = desktopWorkingArea.Right - this.Width;
            this.Top = desktopWorkingArea.Bottom - this.Height;

            WindowState = WindowState.Normal;
        }

        private async void Translate_Clicked(object sender, RoutedEventArgs e)
        {
            var originalText = OriginalText.Text;
            var translatedText = await GetTranslatedText(originalText);
            if(translatedText != null)
            {
                TranslatedText.Text = translatedText.Text;
            }
        }

        private async Task<TranslationResponse> GetTranslatedText(string originalText)
        {
            TranslationResponse result;
            try
            {
                var request = new TranslationRequest
                {
                    Text = originalText,
                    LanguageCode = translateToSwedish ? "sv" : "en"
                };

                using var client = new HttpClient();
                var response = await client.PostAsync(TranslationServiceUrl, new StringContent(JsonConvert.SerializeObject(request), Encoding.UTF8, "application/json"));
                var jsonString = await response.Content.ReadAsStringAsync();
                return JsonConvert.DeserializeObject<TranslationResponse>(jsonString);
            }
            catch(Exception e)
            {
                Console.WriteLine("Error while retrieving translation", e);
            }

            return null;
        }

        private void Swap_Languages_Clicked(object sender, RoutedEventArgs e)
        {
            translateToSwedish = !translateToSwedish;
            UpdateLanguageLabel();
        }

        private void UpdateLanguageLabel()
        {
            if(translateToSwedish)
            {
                LanguageLabel.Content = string.Format(LANGUAGE_LABEL_TEMPLATE, "English", "Swedish");
            }
            else
            {
                LanguageLabel.Content = string.Format(LANGUAGE_LABEL_TEMPLATE, "Swedish", "English");
            }

            var originalText = OriginalText.Text;
            var translatedText = TranslatedText.Text;

            TranslatedText.Text = originalText;
            OriginalText.Text = translatedText;
        }

        private void Toggle_Top_Click(object sender, RoutedEventArgs e)
        {
            Topmost = !Topmost;
        }
    }
}
