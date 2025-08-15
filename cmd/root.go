package cmd

import (
	"fmt"
	"os"
	"raindrop/internal/cli"

	"github.com/spf13/cobra"
)

// Execute 整个程序的入口点
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "\nfor more information, try '--help'\n")
		os.Exit(1)
	}
}

// newRootCmd 私有构造函数, 在这里创建根命令 (Root Command) 配置
func newRootCmd() *cobra.Command {
	// 初始化参数结构体
	runner := cli.NewRunner()

	var cmd = &cobra.Command{
		Use:   "rd [shared-files...]",
		Short: "for file and text sharing",

		SilenceUsage: true,
		Args:         cobra.ArbitraryArgs, // 允许任意数量的位置参数

		// RunE 是执行入口函数, 它允许返回 error, 是 cobra 的推荐的实践
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				runner.Config.SharedPaths = []string{"."} // 默认共享当前路径
			} else {
				runner.Config.SharedPaths = args // 将所有参数赋值给共享路径切片
			}

			if err := runner.Validate(); err != nil {
				return err
			}

			if err := runner.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&runner.Config.Message, "message", "m", "", "Message content to send")
	cmd.Flags().StringVarP(&runner.Config.ContentPath, "content", "c", "", "Path of a file whose content will be sent as plain text")
	cmd.Flags().IntVarP(&runner.Config.Port, "port", "p", 1130, "Specify the server port")

	return cmd
}
