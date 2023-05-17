package cmd

import (
	"fmt"
	"os"

	"subtitle-translator-ai/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "sta",
	Short: "sta 是一个基于AI的字幕文件翻译脚本",
	Long: `Subtitle Translator AI
一个基于AI的字幕文件翻译脚本`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Run()
	},
}

var (
	config       string
	engine       string
	jiebaDictDir string
	language     string
	subject      string
	batchSize    int
	endpoints    string
	proxy        string
	openaiKey    string
	inpath       string
	outpath      string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&config, "config", "f", "", "配置文件")
	rootCmd.PersistentFlags().StringVarP(&engine, "engine", "e", "Google", "翻译引擎")
	rootCmd.PersistentFlags().StringVarP(&jiebaDictDir, "jiebaDictDir", "j", "jiebadict", "结巴分词字典文件目录")
	rootCmd.PersistentFlags().StringVarP(&language, "language", "l", "Chinese", "目标语言")
	rootCmd.PersistentFlags().StringVarP(&subject, "subject", "s", "movie", "目标语言")
	rootCmd.PersistentFlags().IntVarP(&batchSize, "batchSize", "b", 10, "每次翻译的字幕行数")
	rootCmd.PersistentFlags().StringVarP(&endpoints, "endpoints", "d", "", "指定语义完整的字幕序号集合(10,25,40)")
	rootCmd.PersistentFlags().StringVarP(&openaiKey, "openaiKey", "k", "", "openai API key")
	rootCmd.PersistentFlags().StringVarP(&proxy, "proxy", "p", "", "代理配置(127.0.0.1:51837)")
	rootCmd.PersistentFlags().StringVarP(&inpath, "inpath", "i", "", "输入字幕文件路径")
	rootCmd.PersistentFlags().StringVarP(&outpath, "outpath", "o", "", "输出字幕文件路径")
	rootCmd.MarkPersistentFlagRequired("inpath")
	rootCmd.MarkPersistentFlagRequired("outpath")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("engine", rootCmd.PersistentFlags().Lookup("engine"))
	viper.BindPFlag("jiebaDictDir", rootCmd.PersistentFlags().Lookup("jiebaDictDir"))
	viper.BindPFlag("language", rootCmd.PersistentFlags().Lookup("language"))
	viper.BindPFlag("subject", rootCmd.PersistentFlags().Lookup("subject"))
	viper.BindPFlag("batchSize", rootCmd.PersistentFlags().Lookup("batchSize"))
	viper.BindPFlag("endpoints", rootCmd.PersistentFlags().Lookup("endpoints"))
	viper.BindPFlag("openaiKey", rootCmd.PersistentFlags().Lookup("openaiKey"))
	viper.BindPFlag("inpath", rootCmd.PersistentFlags().Lookup("inpath"))
	viper.BindPFlag("outpath", rootCmd.PersistentFlags().Lookup("outpath"))
}

func initConfig() {
	configFile := viper.GetString("config")
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Failed to read config file:", err)
			os.Exit(1)
		}
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
