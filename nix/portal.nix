{
  config,
  lib,
  pkgs,
  ...
}: let
  cfg = config.programs.caelestia-extras.portal;
  qtPluginPath = lib.makeSearchPathOutput "lib" "lib/qt-6/plugins" [
    pkgs.qt6Packages.qt6ct
    pkgs.kdePackages.breeze
  ];
in {
  config = lib.mkIf cfg.enable {
    xdg.configFile = {
      "portal-qt/qt6ct/qt6ct.conf".text = ''
        [Appearance]
        color_scheme_path=${cfg.configHome}/portal-qt/qt6ct/colors/caelestia.conf
        custom_palette=true
        icon_theme=${cfg.iconTheme}
        standard_dialogs=default
        style=Breeze

        [Fonts]
        fixed="Monospace,12,-1,5,50,0,0,0,0,0"
        general="Sans Serif,12,-1,5,50,0,0,0,0,0"

        [Interface]
        activate_item_on_single_click=0
        buttonbox_layout=0
        cursor_flash_time=1000
        dialog_buttons_have_icons=1
        double_click_interval=400
        gui_effects=@Invalid()
        keyboard_scheme=2
        menus_have_icons=true
        show_shortcuts_in_context_menus=true
        toolbutton_style=4
        underline_shortcut=1
        wheel_scroll_lines=3
      '';

      "systemd/user/xdg-desktop-portal-gtk.service.d/10-caelestia-theme.conf".text = ''
        [Service]
        Environment=GTK_THEME=${cfg.themeName}
      '';

      "systemd/user/xdg-desktop-portal-hyprland.service.d/10-caelestia-theme.conf".text = ''
        [Service]
        Environment=QT_QPA_PLATFORMTHEME=qt6ct
        Environment=QT_STYLE_OVERRIDE=Breeze
        Environment=QT_PLUGIN_PATH=${qtPluginPath}
        Environment=XDG_CONFIG_HOME=${cfg.configHome}/portal-qt
      '';

      "caelestia/templates/gtk-portal.css".source = ../assets/manual/templates/gtk-portal.css;
    };

    xdg.dataFile = {
      "themes/${cfg.themeName}/gtk-3.0/base-dark.css".source = ../assets/manual/theme/base-dark.css;
      "themes/${cfg.themeName}/gtk-3.0/base-light.css".source = ../assets/manual/theme/base-light.css;
      "themes/${cfg.themeName}/gtk-4.0/base-dark.css".source = ../assets/manual/theme/base-dark.css;
      "themes/${cfg.themeName}/gtk-4.0/base-light.css".source = ../assets/manual/theme/base-light.css;
    };
  };
}
