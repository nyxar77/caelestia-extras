{
  config,
  lib,
  ...
}: let
  cfg = config.programs.caelestia-extras.portal;
in {
  config = lib.mkIf cfg.enable {

  xdg.configFile."portal-qt/qt6ct/qt6ct.conf".text = ''
    [Appearance]
    color_scheme_path=${cfg.configHome}/portal-qt/qt6ct/colors/caelestia.conf
    custom_palette=true
    icon_theme=${cfg.iconTheme}
    standard_dialogs=default
    style=Fusion

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
    stylesheets=${cfg.configHome}/portal-qt/qt6ct/qss/caelestia.qss
    toolbutton_style=4
    underline_shortcut=1
    wheel_scroll_lines=3
  '';

  xdg.configFile."systemd/user/xdg-desktop-portal-gtk.service.d/10-caelestia-theme.conf".text = ''
    [Service]
    Environment=GTK_THEME=${cfg.themeName}
  '';

  xdg.configFile."systemd/user/xdg-desktop-portal-hyprland.service.d/10-caelestia-theme.conf".text = ''
    [Service]
    Environment=QT_QPA_PLATFORMTHEME=qt6ct
    Environment=QT_STYLE_OVERRIDE=Fusion
    Environment=XDG_CONFIG_HOME=${cfg.configHome}/portal-qt
  '';

  xdg.dataFile."themes/${cfg.themeName}/gtk-3.0/base-dark.css".text = ''
    @import url("resource:///org/gtk/libgtk/theme/Adwaita/gtk-contained-dark.css");
  '';

  xdg.dataFile."themes/${cfg.themeName}/gtk-3.0/base-light.css".text = ''
    @import url("resource:///org/gtk/libgtk/theme/Adwaita/gtk-contained.css");
  '';

  xdg.dataFile."themes/${cfg.themeName}/gtk-4.0/base-dark.css".text = ''
    @import url("resource:///org/gtk/libgtk/theme/Adwaita/gtk-contained-dark.css");
  '';

  xdg.dataFile."themes/${cfg.themeName}/gtk-4.0/base-light.css".text = ''
    @import url("resource:///org/gtk/libgtk/theme/Adwaita/gtk-contained.css");
  '';

  xdg.configFile."caelestia/templates/gtk-portal.css".text = ''
    /*
      XDG portal file chooser. Adwaita supplies the widget behavior; this keeps
      the layout compact and uses the scheme only for contrast and state.
    */

    @import url("base-{{ mode }}.css");

    @define-color portal_accent #{{ primary.hex }};
    @define-color portal_accent_fg #{{ onPrimary.hex }};
    @define-color portal_window #{{ background.hex }};
    @define-color portal_bg #{{ surface.hex }};
    @define-color portal_fg #{{ onSurface.hex }};
    @define-color portal_view #{{ surface.hex }};
    @define-color portal_sidebar #{{ surfaceContainerLow.hex }};
    @define-color portal_panel #{{ surfaceContainer.hex }};
    @define-color portal_panel_high #{{ surfaceContainerHigh.hex }};
    @define-color portal_muted #{{ onSurfaceVariant.hex }};
    @define-color portal_border #{{ outlineVariant.hex }};

    window.background,
    window {
      color: @portal_fg;
      background-color: @portal_window;
      -gtk-icon-theme: "${cfg.iconTheme}";
    }

    filechooser {
      color: @portal_fg;
      background-color: @portal_bg;
    }

    /*
      The chooser is intentionally quieter than the shell: it uses the scheme
      for contrast and state, not as decoration.  Keep this block last so it
      also normalises the slightly different GTK 3 and GTK 4 picker trees.
    */
    window.csd {
      border-radius: 10px;
      box-shadow: none;
    }

    headerbar,
    .titlebar,
    headerbar:focus,
    .titlebar:focus,
    headerbar:backdrop,
    .titlebar:backdrop {
      min-height: 42px;
      padding: 0 7px;
      color: @portal_fg;
      background-color: @portal_panel;
      background-image: none;
      border-top: 0;
      border-bottom: 1px solid @portal_border;
      box-shadow: none;
    }

    headerbar button,
    headerbar button.flat,
    .titlebar button,
    .titlebar button.flat,
    headerbar button:backdrop,
    .titlebar button:backdrop {
      min-height: 28px;
      min-width: 28px;
      margin: 3px 2px;
      padding: 0 9px;
      color: @portal_fg;
      background-color: transparent;
      background-image: none;
      border: 1px solid transparent;
      border-radius: 7px;
      box-shadow: none;
    }

    headerbar button.text-button:not(.default),
    .titlebar button.text-button:not(.default) {
      min-width: 58px;
      padding-left: 11px;
      padding-right: 11px;
    }

    headerbar button:hover,
    headerbar button.flat:hover,
    .titlebar button:hover,
    .titlebar button.flat:hover {
      color: @portal_fg;
      background-color: @portal_panel_high;
      border-color: alpha(@portal_border, 0.84);
      box-shadow: none;
    }

    headerbar button.default,
    headerbar button.suggested-action,
    .titlebar button.default,
    .titlebar button.suggested-action,
    filechooser button.default,
    filechooser button.suggested-action {
      min-height: 28px;
      min-width: 64px;
      padding-left: 12px;
      padding-right: 12px;
      color: @portal_accent_fg;
      background-color: @portal_accent;
      background-image: none;
      border-color: @portal_accent;
      border-radius: 7px;
      box-shadow: none;
    }

    headerbar button.default:hover,
    headerbar button.suggested-action:hover,
    .titlebar button.default:hover,
    .titlebar button.suggested-action:hover,
    filechooser button.default:hover,
    filechooser button.suggested-action:hover {
      color: @portal_accent_fg;
      background-color: shade(@portal_accent, 1.06);
      border-color: shade(@portal_accent, 1.06);
    }

    headerbar button:disabled,
    headerbar button.default:disabled,
    headerbar button.suggested-action:disabled,
    .titlebar button:disabled,
    .titlebar button.default:disabled,
    .titlebar button.suggested-action:disabled,
    filechooser button:disabled,
    filechooser button.default:disabled,
    filechooser button.suggested-action:disabled {
      color: alpha(@portal_muted, 0.72);
      background-color: @portal_panel;
      background-image: none;
      border-color: alpha(@portal_border, 0.72);
      box-shadow: none;
    }

    filechooser #pathbarbox,
    filechooser #pathbarbox:backdrop {
      padding: 2px 5px;
      background-color: @portal_bg;
      border-bottom: 1px solid @portal_border;
    }

    filechooser pathbar button,
    filechooser .path-bar button,
    filechooser pathbar > button,
    filechooser .path-bar > button,
    filechooser #pathbarbox > button,
    filechooser pathbar button:backdrop,
    filechooser .path-bar button:backdrop {
      min-height: 26px;
      min-width: 26px;
      margin: 1px;
      padding: 1px 7px;
      color: @portal_fg;
      background-color: transparent;
      background-image: none;
      border: 1px solid transparent;
      border-radius: 6px;
      box-shadow: none;
    }

    filechooser pathbar button:hover,
    filechooser .path-bar button:hover,
    filechooser #pathbarbox > button:hover {
      color: @portal_fg;
      background-color: @portal_panel;
      border-color: alpha(@portal_border, 0.82);
    }

    filechooser pathbar button:checked,
    filechooser .path-bar button:checked,
    filechooser pathbar > button:checked,
    filechooser .path-bar > button:checked,
    filechooser #pathbarbox > button:checked,
    filechooser pathbar button:backdrop:checked,
    filechooser .path-bar button:backdrop:checked {
      color: @portal_fg;
      background-color: @portal_panel_high;
      background-image: none;
      border-color: alpha(@portal_accent, 0.58);
      box-shadow: none;
    }

    filechooser placessidebar.sidebar,
    filechooser placessidebar,
    filechooser placessidebar.sidebar:backdrop,
    filechooser placessidebar:backdrop {
      min-width: 164px;
      padding: 7px 5px;
      color: @portal_fg;
      background-color: @portal_sidebar;
      background-image: none;
      border-right: 1px solid @portal_border;
    }

    filechooser placessidebar row,
    filechooser placessidebar row:backdrop {
      min-height: 34px;
      margin: 2px 4px;
      padding: 5px 9px;
      color: @portal_muted;
      background-color: transparent;
      background-image: none;
      border: 1px solid transparent;
      border-radius: 7px;
      box-shadow: none;
    }

    filechooser placessidebar row:hover {
      color: @portal_fg;
      background-color: @portal_panel;
      border-color: alpha(@portal_border, 0.62);
      box-shadow: none;
    }

    filechooser placessidebar row:selected,
    filechooser placessidebar row:backdrop:selected {
      color: @portal_fg;
      background-color: alpha(@portal_accent, 0.15);
      background-image: none;
      border-color: alpha(@portal_accent, 0.40);
      box-shadow: none;
    }

    filechooser treeview.view header button,
    filechooser treeview.view header button:backdrop {
      min-height: 32px;
      padding: 0 10px;
      color: @portal_fg;
      background-color: @portal_panel;
      background-image: none;
      border-top: 0;
      border-bottom: 1px solid @portal_border;
      border-left: 0;
      border-right: 1px solid @portal_border;
      border-radius: 0;
      box-shadow: none;
    }

    filechooser treeview.view,
    filechooser treeview.view:backdrop {
      min-height: 32px;
      padding: 4px 10px;
      color: @portal_fg;
      background-color: @portal_view;
      background-image: none;
      box-shadow: none;
    }

    filechooser treeview.view:hover {
      background-color: alpha(@portal_panel, 0.76);
      background-image: none;
    }

    filechooser treeview.view:selected,
    filechooser treeview.view:selected:focus,
    filechooser treeview.view:selected:backdrop {
      color: @portal_fg;
      background-color: alpha(@portal_accent, 0.13);
      background-image: none;
      border-color: transparent;
      box-shadow: none;
    }

    filechooser .dialog-action-box,
    filechooser > box > actionbar,
    filechooser placesview > actionbar > revealer > box,
    filechooser scrolledwindow + actionbar > revealer > box,
    filechooser .dialog-action-box:backdrop,
    filechooser actionbar:backdrop {
      min-height: 36px;
      padding: 3px 7px;
      color: @portal_fg;
      background-color: @portal_bg;
      border-top: 1px solid @portal_border;
      box-shadow: none;
    }

    filechooser checkbutton {
      min-height: 26px;
      padding: 1px 3px;
      color: @portal_muted;
    }

    filechooser checkbutton check,
    filechooser checkbutton > check,
    filechooser actionbar checkbutton check,
    filechooser .dialog-action-box checkbutton check {
      min-width: 13px;
      min-height: 13px;
      margin-right: 5px;
      background-color: @portal_panel;
      background-image: none;
      border: 1px solid @portal_border;
      border-radius: 3px;
      box-shadow: none;
    }

    filechooser checkbutton check:checked,
    filechooser checkbutton > check:checked,
    filechooser actionbar checkbutton check:checked,
    filechooser .dialog-action-box checkbutton check:checked {
      color: @portal_accent_fg;
      background-color: @portal_accent;
      border-color: @portal_accent;
    }

    filechooser scrollbar slider,
    filechooser scrollbar slider:backdrop {
      min-width: 7px;
      min-height: 7px;
      border: 2px solid transparent;
      background-color: alpha(@portal_muted, 0.46);
      background-image: none;
      box-shadow: none;
    }

    filechooser scrollbar slider:hover {
      background-color: alpha(@portal_muted, 0.72);
    }
  '';

  xdg.configFile."caelestia/templates/gtk-global.css".text = ''
    /*
      Dynamic GTK colours only. Adwaita/libadwaita remains responsible for
      widget layout, typography, spacing, animations, and app-specific styles.
    */

    @define-color accent_color #{{ primary.hex }};
    @define-color accent_bg_color #{{ primary.hex }};
    @define-color accent_fg_color #{{ onPrimary.hex }};
    @define-color destructive_color #{{ error.hex }};
    @define-color destructive_bg_color #{{ error.hex }};
    @define-color destructive_fg_color #{{ onError.hex }};
    @define-color success_color #{{ success.hex }};
    @define-color success_bg_color #{{ successContainer.hex }};
    @define-color success_fg_color #{{ onSuccessContainer.hex }};
    @define-color warning_color #{{ secondary.hex }};
    @define-color warning_bg_color #{{ secondaryContainer.hex }};
    @define-color warning_fg_color #{{ onSecondaryContainer.hex }};
    @define-color error_color #{{ error.hex }};
    @define-color error_bg_color #{{ errorContainer.hex }};
    @define-color error_fg_color #{{ onErrorContainer.hex }};

    @define-color window_bg_color #{{ surface.hex }};
    @define-color window_fg_color #{{ onSurface.hex }};
    @define-color view_bg_color #{{ surface.hex }};
    @define-color view_fg_color #{{ onSurface.hex }};
    @define-color headerbar_bg_color #{{ surfaceContainerLow.hex }};
    @define-color headerbar_fg_color #{{ onSurface.hex }};
    @define-color headerbar_border_color #{{ outlineVariant.hex }};
    @define-color sidebar_bg_color #{{ surfaceContainerLow.hex }};
    @define-color sidebar_fg_color #{{ onSurface.hex }};
    @define-color sidebar_backdrop_color #{{ surface.hex }};
    @define-color sidebar_border_color #{{ outlineVariant.hex }};
    @define-color card_bg_color #{{ surfaceContainer.hex }};
    @define-color card_fg_color #{{ onSurface.hex }};
    @define-color dialog_bg_color #{{ surfaceContainer.hex }};
    @define-color dialog_fg_color #{{ onSurface.hex }};
    @define-color popover_bg_color #{{ surfaceContainerHigh.hex }};
    @define-color popover_fg_color #{{ onSurface.hex }};
    @define-color shade_color alpha(#{{ shadow.hex }}, 0.32);
    @define-color scrollbar_outline_color alpha(#{{ outline.hex }}, 0.38);
    @define-color borders #{{ outlineVariant.hex }};
    @define-color unfocused_borders alpha(#{{ outlineVariant.hex }}, 0.72);

    @define-color theme_fg_color #{{ onSurface.hex }};
    @define-color theme_text_color #{{ onSurface.hex }};
    @define-color theme_bg_color #{{ surface.hex }};
    @define-color theme_base_color #{{ surface.hex }};
    @define-color theme_selected_bg_color alpha(#{{ primary.hex }}, 0.32);
    @define-color theme_selected_fg_color #{{ onSurface.hex }};
    @define-color insensitive_fg_color alpha(#{{ onSurfaceVariant.hex }}, 0.62);
    @define-color insensitive_bg_color #{{ surfaceContainerLow.hex }};
    @define-color insensitive_base_color #{{ surfaceContainerLowest.hex }};
    @define-color theme_unfocused_fg_color #{{ onSurfaceVariant.hex }};
    @define-color theme_unfocused_text_color #{{ onSurfaceVariant.hex }};
    @define-color theme_unfocused_bg_color #{{ surfaceDim.hex }};
    @define-color theme_unfocused_base_color #{{ surfaceDim.hex }};
    @define-color theme_unfocused_selected_bg_color alpha(#{{ secondary.hex }}, 0.28);
    @define-color theme_unfocused_selected_fg_color #{{ onSurface.hex }};

    /*
      Keep GTK's native geometry and interaction model, but make the active
      Caelestia palette visible in ordinary GTK applications.
    */
    window,
    dialog,
    .background {
      color: @window_fg_color;
      background-color: @window_bg_color;
    }

    headerbar,
    .titlebar {
      color: @headerbar_fg_color;
      background-color: @headerbar_bg_color;
      background-image: linear-gradient(to bottom, alpha(@accent_color, 0.10), transparent 48%);
      border-bottom-color: @headerbar_border_color;
    }

    .sidebar,
    placessidebar,
    .navigation-sidebar {
      color: @sidebar_fg_color;
      background-color: @sidebar_bg_color;
      border-color: @sidebar_border_color;
    }

    view,
    textview,
    treeview,
    list {
      color: @view_fg_color;
      background-color: @view_bg_color;
    }

    button,
    entry,
    combobox button,
    spinbutton {
      color: @window_fg_color;
      background-color: @card_bg_color;
      border-color: @borders;
    }

    button:hover,
    button:focus-visible,
    entry:focus,
    combobox button:hover,
    spinbutton:focus {
      background-color: @popover_bg_color;
      border-color: @accent_color;
    }

    button:checked,
    button:active,
    row:selected,
    treeview.view:selected,
    textview text selection,
    entry selection {
      color: @theme_selected_fg_color;
      background-color: @theme_selected_bg_color;
    }

    row:hover,
    treeview.view:hover {
      background-color: alpha(@accent_color, 0.10);
    }

  '';

  xdg.configFile."caelestia/templates/gtk.css".text = ''
    @define-color accent_color #{{ primary.hex }};
    @define-color accent_fg_color #{{ onPrimary.hex }};
    @define-color accent_bg_color #{{ primary.hex }};
    @define-color primary_color #{{ primary.hex }};
    @define-color primary_fg_color #{{ onPrimary.hex }};
    @define-color secondary_color #{{ secondary.hex }};
    @define-color secondary_fg_color #{{ onSecondary.hex }};
    @define-color tertiary_color #{{ tertiary.hex }};
    @define-color tertiary_fg_color #{{ onTertiary.hex }};
    @define-color window_bg_color #{{ surface.hex }};
    @define-color window_fg_color #{{ onSurface.hex }};
    @define-color headerbar_bg_color #{{ surfaceContainerLow.hex }};
    @define-color headerbar_fg_color #{{ onSurface.hex }};
    @define-color popover_bg_color #{{ surfaceContainer.hex }};
    @define-color popover_fg_color #{{ onSurface.hex }};
    @define-color view_bg_color #{{ surface.hex }};
    @define-color view_fg_color #{{ onSurface.hex }};
    @define-color card_bg_color #{{ surfaceContainer.hex }};
    @define-color card_fg_color #{{ onSurface.hex }};
    @define-color sidebar_bg_color #{{ surfaceContainerLow.hex }};
    @define-color sidebar_fg_color #{{ onSurface.hex }};
    @define-color sidebar_border_color #{{ outlineVariant.hex }};
    @define-color sidebar_backdrop_color #{{ surface.hex }};
    @define-color border_color #{{ outlineVariant.hex }};
    @define-color muted_color #{{ onSurfaceVariant.hex }};
    @define-color hover_color #{{ surfaceContainerHigh.hex }};
    @define-color active_color #{{ primary.hex }};
    @define-color active_fg_color #{{ onPrimary.hex }};
    @define-color focus_ring alpha(@accent_color, 0.55);
    @define-color theme_selected_bg_color alpha(@accent_color, 0.32);
    @define-color theme_selected_fg_color @accent_fg_color;
    @define-color layer_0 #{{ surface.hex }};
    @define-color layer_1 #{{ surfaceContainerLow.hex }};
    @define-color layer_2 #{{ surfaceContainer.hex }};
    @define-color layer_3 #{{ surfaceContainerHigh.hex }};
    @define-color layer_4 #{{ surfaceContainerHighest.hex }};
    @define-color accent_wash alpha(@accent_color, 0.14);
    @define-color accent_wash_strong alpha(@accent_color, 0.24);

    * {
      outline-color: alpha(@accent_color, 0.42);
      -gtk-icon-shadow: none;
    }

    window,
    dialog,
    messagedialog,
    filechooser,
    popover,
    menu,
    .background {
      color: @window_fg_color;
      background-color: @window_bg_color;
    }

    .csd,
    decoration {
      border-radius: 9px;
      box-shadow: 0 18px 46px alpha(black, 0.42);
    }

    headerbar,
    .titlebar {
      min-height: 46px;
      color: @headerbar_fg_color;
      background-color: @layer_1;
      background-image: linear-gradient(to bottom, alpha(@accent_color, 0.12), transparent 42%);
      border-color: @border_color;
      border-top: 2px solid alpha(@accent_color, 0.72);
      border-bottom: 1px solid @border_color;
      box-shadow: inset 0 -1px alpha(@accent_color, 0.10), 0 1px 0 alpha(white, 0.03);
    }

    headerbar:backdrop,
    .titlebar:backdrop,
    window:backdrop headerbar,
    window:backdrop .titlebar {
      color: @muted_color;
      background-color: @layer_1;
      background-image: none;
      border-top-color: alpha(@accent_color, 0.34);
      border-bottom-color: @border_color;
      box-shadow: none;
    }

    headerbar:focus,
    headerbar:focus-within,
    .titlebar:focus,
    .titlebar:focus-within,
    window:focus headerbar,
    window:focus .titlebar {
      color: @headerbar_fg_color;
      background-color: @layer_1;
      background-image: linear-gradient(to bottom, alpha(@accent_color, 0.18), transparent 48%);
      border-top-color: @accent_color;
      border-bottom-color: alpha(@accent_color, 0.42);
    }

    headerbar *,
    .titlebar * {
      background-image: none;
    }

    headerbar button,
    headerbar menubutton,
    headerbar combobox button,
    .titlebar button,
    .titlebar menubutton,
    .titlebar combobox button {
      min-height: 30px;
      min-width: 30px;
      padding: 4px 9px;
      color: @headerbar_fg_color;
      background-color: alpha(@layer_3, 0.56);
      border-color: alpha(@border_color, 0.78);
      border-radius: 8px;
      box-shadow: inset 0 1px alpha(white, 0.05), 0 1px 4px alpha(black, 0.20);
    }

    headerbar button:hover,
    headerbar menubutton:hover,
    headerbar combobox button:hover,
    .titlebar button:hover,
    .titlebar menubutton:hover,
    .titlebar combobox button:hover {
      color: @headerbar_fg_color;
      background-color: @accent_wash_strong;
      border-color: @accent_color;
      box-shadow: inset 0 0 0 1px alpha(@accent_color, 0.22), 0 2px 8px alpha(black, 0.26);
    }

    popover,
    menu,
    .csd popover {
      border: 1px solid @border_color;
      border-radius: 8px;
      background-color: @popover_bg_color;
      box-shadow: 0 12px 34px alpha(black, 0.34);
    }

    button,
    menubutton,
    combobox button {
      min-height: 32px;
      padding: 6px 11px;
      color: @window_fg_color;
      background-color: @layer_2;
      background-image: linear-gradient(to bottom, alpha(white, 0.035), alpha(black, 0.035));
      border: 1px solid @border_color;
      border-radius: 7px;
      box-shadow: inset 0 1px alpha(white, 0.04), 0 1px 2px alpha(black, 0.18);
    }

    button:hover,
    menubutton:hover,
    combobox button:hover {
      color: @window_fg_color;
      background-color: @layer_3;
      background-image: linear-gradient(to bottom, @accent_wash, transparent);
      border-color: @accent_color;
      box-shadow: inset 0 0 0 1px alpha(@accent_color, 0.22), 0 2px 5px alpha(black, 0.22);
    }

    button:checked,
    button:active,
    menubutton:active {
      color: @active_fg_color;
      background: @active_color;
      border-color: @accent_color;
      box-shadow: none;
    }

    button.suggested-action,
    button.default {
      color: @accent_fg_color;
      background: @accent_color;
      border-color: @accent_color;
    }

    button.destructive-action {
      color: @accent_fg_color;
      background: #{{ error.hex }};
      border-color: #{{ error.hex }};
    }

    entry,
    searchentry,
    spinbutton,
    textview,
    textview text,
    treeview,
    list,
    listview,
    columnview,
    gridview,
    placessidebar,
    scrolledwindow {
      color: @view_fg_color;
      background-color: @view_bg_color;
      border-color: @border_color;
    }

    scrolledwindow,
    viewport,
    notebook,
    stack,
    paned {
      background-color: @view_bg_color;
    }

    entry,
    searchentry,
    spinbutton {
      min-height: 34px;
      padding: 6px 10px;
      color: @view_fg_color;
      background-color: @layer_1;
      background-image: none;
      border: 1px solid @border_color;
      border-radius: 7px;
      box-shadow: inset 0 0 0 1px alpha(black, 0.10);
    }

    entry:focus,
    searchentry:focus,
    spinbutton:focus {
      border-color: @accent_color;
      box-shadow: inset 0 0 0 1px @focus_ring, 0 0 0 3px alpha(@accent_color, 0.10);
    }

    placessidebar,
    placessidebar list,
    filechooser .sidebar,
    .sidebar {
      color: @sidebar_fg_color;
      background-color: @layer_1;
      border-right: 1px solid @border_color;
    }

    placessidebar row,
    .sidebar row {
      min-height: 34px;
      margin: 3px 6px;
      padding: 5px 8px;
      color: @sidebar_fg_color;
      background-color: transparent;
      border-radius: 8px;
      box-shadow: none;
    }

    placessidebar row:hover,
    .sidebar row:hover {
      background-color: @layer_3;
    }

    placessidebar row:selected,
    .sidebar row:selected {
      color: @window_fg_color;
      background-color: @accent_wash_strong;
      box-shadow: inset 3px 0 @accent_color;
    }

    toolbar,
    actionbar,
    searchbar {
      color: @window_fg_color;
      background-color: @headerbar_bg_color;
      border-color: @border_color;
    }

    row,
    list row,
    treeview.view,
    columnview row,
    gridview child {
      min-height: 30px;
      color: @view_fg_color;
      background-color: transparent;
      border-radius: 7px;
      border: 1px solid transparent;
    }

    row:hover,
    list row:hover,
    columnview row:hover,
    gridview child:hover {
      background-color: @layer_3;
      border-color: alpha(@accent_color, 0.34);
    }

    row:selected,
    list row:selected,
    treeview.view:selected,
    columnview row:selected,
    gridview child:selected {
      color: @window_fg_color;
      background-color: @accent_wash_strong;
      border-color: alpha(@accent_color, 0.56);
      box-shadow: inset 3px 0 @accent_color;
    }

    selection,
    text selection,
    entry selection {
      color: @active_fg_color;
      background-color: @accent_color;
    }

    scrollbar,
    scrollbar trough {
      background-color: transparent;
      border: none;
    }

    scrollbar slider {
      min-width: 6px;
      min-height: 6px;
      border: 2px solid transparent;
      border-radius: 999px;
      background-color: alpha(@muted_color, 0.42);
    }

    scrollbar slider:hover {
      background-color: alpha(@accent_color, 0.72);
    }

    tabs,
    tab {
      color: @muted_color;
      background-color: transparent;
    }

    tab:hover {
      color: @window_fg_color;
      background-color: @hover_color;
    }

    tab:checked {
      color: @window_fg_color;
      background-color: @card_bg_color;
      box-shadow: inset 0 -2px @accent_color;
    }

    separator {
      background-color: @border_color;
    }

    label,
    image,
    .dim-label {
      color: inherit;
    }

    .dim-label,
    label:disabled,
    button:disabled {
      color: @muted_color;
    }

    @import "thunar.css";

    window,
    dialog,
    filechooser,
    .background {
      background-image: linear-gradient(135deg, alpha(@secondary_color, 0.08), transparent 28%, alpha(@tertiary_color, 0.07));
    }

    headerbar,
    .titlebar,
    headerbar:focus,
    headerbar:focus-within,
    .titlebar:focus,
    .titlebar:focus-within,
    window:focus headerbar,
    window:focus .titlebar,
    window:backdrop headerbar,
    window:backdrop .titlebar {
      color: @window_fg_color;
      background-color: @layer_1;
      background-image: linear-gradient(120deg, alpha(@primary_color, 0.24), alpha(@tertiary_color, 0.14) 48%, alpha(@secondary_color, 0.16));
      border-top: 2px solid @primary_color;
      border-bottom: 1px solid alpha(@tertiary_color, 0.42);
      box-shadow: inset 0 -1px alpha(@primary_color, 0.20), 0 1px 0 alpha(white, 0.04);
    }

    headerbar:backdrop,
    .titlebar:backdrop,
    window:backdrop headerbar,
    window:backdrop .titlebar {
      background-image: linear-gradient(120deg, alpha(@primary_color, 0.14), alpha(@tertiary_color, 0.08));
      border-top-color: alpha(@primary_color, 0.52);
      border-bottom-color: @border_color;
    }

    headerbar button,
    .titlebar button,
    headerbar menubutton,
    .titlebar menubutton {
      background-image: linear-gradient(135deg, alpha(@secondary_color, 0.16), alpha(@tertiary_color, 0.10));
      border-color: alpha(@tertiary_color, 0.42);
    }

    headerbar button:hover,
    .titlebar button:hover,
    headerbar menubutton:hover,
    .titlebar menubutton:hover {
      background-image: linear-gradient(135deg, alpha(@primary_color, 0.30), alpha(@tertiary_color, 0.24));
      border-color: @primary_color;
    }

    button,
    menubutton,
    combobox button {
      background-image: linear-gradient(145deg, alpha(@secondary_color, 0.13), alpha(@layer_3, 0.94));
      border-color: alpha(@secondary_color, 0.36);
    }

    button:hover,
    menubutton:hover,
    combobox button:hover {
      background-image: linear-gradient(135deg, alpha(@primary_color, 0.22), alpha(@tertiary_color, 0.20));
      border-color: alpha(@primary_color, 0.78);
    }

    button:checked,
    button:active,
    button.suggested-action,
    button.default {
      color: @primary_fg_color;
      background-image: linear-gradient(135deg, @primary_color, @tertiary_color);
      border-color: @primary_color;
    }

    entry,
    searchentry,
    spinbutton {
      background-image: linear-gradient(135deg, alpha(@secondary_color, 0.10), alpha(@layer_1, 0.96));
      border-color: alpha(@secondary_color, 0.34);
    }

    entry:focus,
    searchentry:focus,
    spinbutton:focus {
      background-image: linear-gradient(135deg, alpha(@primary_color, 0.12), alpha(@tertiary_color, 0.10));
      border-color: @primary_color;
      box-shadow: inset 0 0 0 1px alpha(@primary_color, 0.48), 0 0 0 3px alpha(@tertiary_color, 0.16);
    }

    placessidebar,
    placessidebar list,
    filechooser .sidebar,
    .sidebar {
      background-image: linear-gradient(180deg, alpha(@secondary_color, 0.12), alpha(@layer_1, 0.96) 38%, alpha(@tertiary_color, 0.08));
      border-right-color: alpha(@secondary_color, 0.38);
    }

    placessidebar row:hover,
    .sidebar row:hover {
      background-image: linear-gradient(90deg, alpha(@secondary_color, 0.24), alpha(@tertiary_color, 0.12));
      box-shadow: inset 3px 0 alpha(@secondary_color, 0.72);
    }

    placessidebar row:selected,
    .sidebar row:selected {
      color: @window_fg_color;
      background-image: linear-gradient(90deg, alpha(@secondary_color, 0.36), alpha(@primary_color, 0.18));
      border-color: alpha(@secondary_color, 0.58);
      box-shadow: inset 4px 0 @secondary_color;
    }

    row:hover,
    list row:hover,
    treeview.view:hover,
    columnview row:hover,
    gridview child:hover {
      background-image: linear-gradient(90deg, alpha(@tertiary_color, 0.16), alpha(@secondary_color, 0.10));
      border-color: alpha(@tertiary_color, 0.44);
    }

    row:selected,
    list row:selected,
    treeview.view:selected,
    columnview row:selected,
    gridview child:selected {
      color: @window_fg_color;
      background-image: linear-gradient(90deg, alpha(@primary_color, 0.34), alpha(@tertiary_color, 0.22));
      border-color: alpha(@primary_color, 0.64);
      box-shadow: inset 4px 0 @primary_color;
    }

    tabs,
    tab {
      background-image: none;
    }

    tab:hover {
      background-image: linear-gradient(135deg, alpha(@tertiary_color, 0.18), alpha(@secondary_color, 0.10));
    }

    tab:checked {
      background-image: linear-gradient(135deg, alpha(@primary_color, 0.20), alpha(@tertiary_color, 0.16));
      box-shadow: inset 0 -3px @tertiary_color;
    }

    scrollbar slider {
      background-image: linear-gradient(to bottom, @secondary_color, @tertiary_color);
    }

    scrollbar slider:hover {
      background-image: linear-gradient(to bottom, @primary_color, @tertiary_color);
    }
  '';

  xdg.configFile."caelestia/templates/qbittorrent.qss".text = ''
    QWidget { background-color: #{{ surface.hex }}; color: #{{ onSurface.hex }}; }
    QMainWindow, QDialog, QFrame, QScrollArea { background-color: #{{ surface.hex }}; }
    QToolBar, QStatusBar, QMenuBar { background-color: #{{ surfaceContainerLow.hex }}; border: 0; border-bottom: 1px solid #{{ outlineVariant.hex }}; }
    QTabWidget::pane { background-color: #{{ surface.hex }}; border: 1px solid #{{ outlineVariant.hex }}; }
    QTabBar::tab { padding: 7px 14px; color: #{{ onSurfaceVariant.hex }}; background-color: #{{ surfaceContainerLow.hex }}; border: 0; border-bottom: 3px solid transparent; }
    QTabBar::tab:hover { color: #{{ onSurface.hex }}; background-color: #{{ surfaceContainer.hex }}; border-bottom-color: #{{ tertiary.hex }}; }
    QTabBar::tab:selected { color: #{{ onSurface.hex }}; background-color: #{{ surfaceContainerHigh.hex }}; border-bottom-color: #{{ primary.hex }}; }
    QPushButton, QToolButton, QComboBox, QLineEdit { min-height: 30px; padding: 4px 10px; background-color: #{{ surfaceContainer.hex }}; border: 1px solid #{{ outlineVariant.hex }}; border-radius: 7px; }
    QPushButton:hover, QToolButton:hover, QComboBox:hover, QLineEdit:focus { border-color: #{{ primary.hex }}; background-color: #{{ surfaceContainerHigh.hex }}; }
    QHeaderView::section { padding: 6px 10px; background-color: #{{ surfaceContainerLow.hex }}; color: #{{ onSurface.hex }}; border: 0; border-right: 1px solid #{{ outlineVariant.hex }}; border-bottom: 1px solid #{{ outlineVariant.hex }}; }
    QTreeView::item { padding: 5px; border: 0; }
    QTreeView::item:hover { background-color: #{{ surfaceContainer.hex }}; }
    QTreeView::item:selected { background-color: #{{ primaryContainer.hex }}; color: #{{ onPrimaryContainer.hex }}; }
    QProgressBar { border: 1px solid #{{ outlineVariant.hex }}; border-radius: 5px; background-color: #{{ surfaceContainerHigh.hex }}; text-align: center; }
    QProgressBar::chunk { background-color: #{{ primary.hex }}; border-radius: 4px; }
  '';

  xdg.configFile."caelestia/templates/qbittorrent.json".text = ''
    { "colors": {
      "Palette.Window": "#{{ surface.hex }}", "Palette.WindowText": "#{{ onSurface.hex }}",
      "Palette.Base": "#{{ surfaceContainerLow.hex }}", "Palette.AlternateBase": "#{{ surfaceContainer.hex }}",
      "Palette.Text": "#{{ onSurface.hex }}", "Palette.ToolTipBase": "#{{ inverseSurface.hex }}", "Palette.ToolTipText": "#{{ inverseOnSurface.hex }}",
      "Palette.Highlight": "#{{ primary.hex }}", "Palette.HighlightedText": "#{{ onPrimary.hex }}",
      "Palette.Button": "#{{ surfaceContainer.hex }}", "Palette.ButtonText": "#{{ onSurface.hex }}",
      "Palette.Link": "#{{ primary.hex }}", "Palette.LinkVisited": "#{{ tertiary.hex }}",
      "Palette.Light": "#{{ surfaceContainerHighest.hex }}", "Palette.Midlight": "#{{ surfaceContainerHigh.hex }}", "Palette.Mid": "#{{ outlineVariant.hex }}", "Palette.Dark": "#{{ surfaceDim.hex }}", "Palette.Shadow": "#{{ shadow.hex }}",
      "Log.Time": "#{{ onSurfaceVariant.hex }}", "Log.Normal": "#{{ onSurface.hex }}", "Log.Info": "#{{ primary.hex }}", "Log.Warning": "#{{ secondary.hex }}", "Log.Critical": "#{{ error.hex }}", "Log.BannedPeer": "#{{ error.hex }}",
      "TransferList.Downloading": "#{{ primary.hex }}", "TransferList.StalledDownloading": "#{{ onSurfaceVariant.hex }}", "TransferList.DownloadingMetadata": "#{{ tertiary.hex }}", "TransferList.ForcedDownloading": "#{{ primary.hex }}", "TransferList.Allocating": "#{{ secondary.hex }}", "TransferList.Uploading": "#{{ tertiary.hex }}", "TransferList.StalledUploading": "#{{ onSurfaceVariant.hex }}", "TransferList.ForcedUploading": "#{{ tertiary.hex }}", "TransferList.QueuedDownloading": "#{{ secondary.hex }}", "TransferList.QueuedUploading": "#{{ secondary.hex }}", "TransferList.CheckingDownloading": "#{{ tertiary.hex }}", "TransferList.CheckingUploading": "#{{ tertiary.hex }}", "TransferList.CheckingResumeData": "#{{ tertiary.hex }}", "TransferList.PausedDownloading": "#{{ onSurfaceVariant.hex }}", "TransferList.PausedUploading": "#{{ onSurfaceVariant.hex }}", "TransferList.Moving": "#{{ secondary.hex }}", "TransferList.MissingFiles": "#{{ error.hex }}", "TransferList.Error": "#{{ error.hex }}"
    } }
  '';

  xdg.configFile."caelestia/templates/pavucontrol-qt.qss".text = ''
    QDialog#MainWindow {
      background-color: #{{ surface.hex }};
      color: #{{ onSurface.hex }};
    }

    QTabWidget#notebook::pane {
      background-color: #{{ surface.hex }};
      border: 0;
      border-top: 1px solid #{{ outlineVariant.hex }};
      top: -1px;
    }

    QTabBar {
      background-color: #{{ surfaceContainerLow.hex }};
    }

    QTabBar::tab {
      min-width: 112px;
      min-height: 34px;
      padding: 7px 14px;
      margin: 0;
      color: #{{ onSurfaceVariant.hex }};
      background-color: #{{ surfaceContainerLow.hex }};
      border: 0;
      border-bottom: 2px solid transparent;
    }

    QTabBar::tab:hover {
      color: #{{ onSurface.hex }};
      background-color: #{{ surfaceContainer.hex }};
      border-bottom-color: #{{ tertiary.hex }};
    }

    QTabBar::tab:selected {
      color: #{{ onSurface.hex }};
      background-color: #{{ surfaceContainerHigh.hex }};
      border-bottom-color: #{{ primary.hex }};
    }

    QScrollArea,
    QScrollArea > QWidget,
    QScrollArea > QWidget > QWidget {
      background-color: #{{ surface.hex }};
      border: 0;
    }

    QWidget#streamsVBox,
    QWidget#recsVBox,
    QWidget#sinksVBox,
    QWidget#sourcesVBox,
    QWidget#cardsVBox {
      background-color: #{{ surface.hex }};
    }

    QWidget#StreamWidget,
    QWidget#DeviceWidget,
    QWidget#CardWidget {
      margin: 6px;
      padding: 12px;
      color: #{{ onSurface.hex }};
      background-color: #{{ surfaceContainerLow.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 10px;
    }

    QWidget#StreamWidget:hover,
    QWidget#DeviceWidget:hover,
    QWidget#CardWidget:hover {
      background-color: #{{ surfaceContainer.hex }};
      border-color: #{{ tertiary.hex }};
    }

    QLabel {
      color: #{{ onSurface.hex }};
      background-color: transparent;
    }

    QLabel#nameLabel,
    QLabel#channelLabel,
    QLabel#volumeLabel,
    QLabel#directionLabel,
    QLabel#noStreamsLabel,
    QLabel#noRecsLabel,
    QLabel#noSinksLabel,
    QLabel#noSourcesLabel,
    QLabel#noCardsLabel {
      color: #{{ onSurfaceVariant.hex }};
    }

    QLabel#boldNameLabel {
      color: #{{ onSurface.hex }};
      font-weight: 600;
    }

    QFrame#line {
      min-height: 1px;
      max-height: 1px;
      margin: 8px 0 0 0;
      color: #{{ outlineVariant.hex }};
      background-color: #{{ outlineVariant.hex }};
      border: 0;
    }

    QToolButton {
      min-width: 32px;
      min-height: 32px;
      padding: 4px;
      color: #{{ onSurface.hex }};
      background-color: #{{ surfaceContainerHigh.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 8px;
    }

    QToolButton:hover {
      background-color: #{{ surfaceContainerHighest.hex }};
      border-color: #{{ tertiary.hex }};
    }

    QToolButton:pressed,
    QToolButton:checked {
      color: #{{ onPrimary.hex }};
      background-color: #{{ primary.hex }};
      border-color: #{{ primary.hex }};
    }

    QToolButton#muteToggleButton:checked {
      color: #{{ onError.hex }};
      background-color: #{{ error.hex }};
      border-color: #{{ error.hex }};
    }

    QToolButton#defaultToggleButton:checked {
      color: #{{ onPrimary.hex }};
      background-color: #{{ primary.hex }};
      border-color: #{{ primary.hex }};
    }

    QComboBox,
    QSpinBox {
      min-height: 32px;
      padding: 5px 30px 5px 10px;
      color: #{{ onSurface.hex }};
      selection-color: #{{ onPrimary.hex }};
      selection-background-color: #{{ primary.hex }};
      background-color: #{{ surfaceContainerHigh.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 8px;
    }

    QComboBox:hover,
    QSpinBox:hover {
      background-color: #{{ surfaceContainerHighest.hex }};
      border-color: #{{ tertiary.hex }};
    }

    QComboBox:focus,
    QSpinBox:focus {
      border-color: #{{ primary.hex }};
    }

    QComboBox QAbstractItemView {
      padding: 5px;
      color: #{{ onSurface.hex }};
      selection-color: #{{ onPrimary.hex }};
      selection-background-color: #{{ primary.hex }};
      background-color: #{{ surfaceContainer.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 8px;
      outline: 0;
    }

    QCheckBox {
      spacing: 8px;
      color: #{{ onSurfaceVariant.hex }};
      background-color: transparent;
    }

    QCheckBox:hover {
      color: #{{ onSurface.hex }};
    }

    QSlider::groove:horizontal {
      height: 6px;
      background-color: #{{ surfaceContainerHighest.hex }};
      border: 0;
      border-radius: 3px;
    }

    QSlider::sub-page:horizontal {
      background-color: #{{ primary.hex }};
      border-radius: 3px;
    }

    QSlider::add-page:horizontal {
      background-color: #{{ surfaceContainerHighest.hex }};
      border-radius: 3px;
    }

    QSlider::handle:horizontal {
      width: 16px;
      margin: -5px 0;
      background-color: #{{ primary.hex }};
      border: 2px solid #{{ onPrimary.hex }};
      border-radius: 8px;
    }

    QSlider::handle:horizontal:hover {
      background-color: #{{ tertiary.hex }};
      border-color: #{{ onTertiary.hex }};
    }

    QProgressBar {
      min-height: 6px;
      max-height: 6px;
      background-color: #{{ surfaceContainerHighest.hex }};
      border: 0;
      border-radius: 3px;
      text-align: center;
    }

    QProgressBar::chunk {
      background-color: #{{ tertiary.hex }};
      border-radius: 3px;
    }

    QScrollBar:vertical {
      width: 10px;
      margin: 2px;
      background-color: transparent;
      border: 0;
    }

    QScrollBar::handle:vertical {
      min-height: 28px;
      background-color: #{{ outlineVariant.hex }};
      border-radius: 5px;
    }

    QScrollBar::handle:vertical:hover {
      background-color: #{{ tertiary.hex }};
    }

    QScrollBar::add-line:vertical,
    QScrollBar::sub-line:vertical,
    QScrollBar::add-page:vertical,
    QScrollBar::sub-page:vertical {
      height: 0;
      background-color: transparent;
      border: 0;
    }

    QToolTip {
      padding: 6px;
      color: #{{ inverseOnSurface.hex }};
      background-color: #{{ inverseSurface.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 6px;
    }
  '';

  xdg.configFile."caelestia/templates/qt6ct-caelestia.conf".text = ''
    [ColorScheme]
    active_colors=#ff{{ onSurface.hex }}, #ff{{ surfaceContainer.hex }}, #ff{{ surfaceContainerHigh.hex }}, #ff{{ surfaceContainerHighest.hex }}, #ff{{ surfaceDim.hex }}, #ff{{ surfaceContainerHigh.hex }}, #ff{{ onSurface.hex }}, #ff{{ onPrimary.hex }}, #ff{{ onSurface.hex }}, #ff{{ surface.hex }}, #ff{{ surfaceContainerLow.hex }}, #ff{{ shadow.hex }}, #ff{{ primary.hex }}, #ff{{ onPrimary.hex }}, #ff{{ secondary.hex }}, #ff{{ tertiary.hex }}, #ff{{ surfaceContainer.hex }}, #ff{{ onSurface.hex }}, #ff{{ inverseSurface.hex }}, #ff{{ inverseOnSurface.hex }}, #99{{ onSurfaceVariant.hex }}
    disabled_colors=#ff{{ onSurfaceVariant.hex }}, #ff{{ surfaceContainer.hex }}, #ff{{ surfaceContainerHigh.hex }}, #ff{{ surfaceContainerHighest.hex }}, #ff{{ surfaceDim.hex }}, #ff{{ surfaceContainerHigh.hex }}, #ff{{ onSurfaceVariant.hex }}, #ff{{ onPrimary.hex }}, #ff{{ onSurfaceVariant.hex }}, #ff{{ surface.hex }}, #ff{{ surfaceContainerLow.hex }}, #ff{{ shadow.hex }}, #ff{{ primary.hex }}, #ff{{ onSurfaceVariant.hex }}, #ff{{ secondary.hex }}, #ff{{ tertiary.hex }}, #ff{{ surfaceContainer.hex }}, #ff{{ onSurface.hex }}, #ff{{ inverseSurface.hex }}, #ff{{ inverseOnSurface.hex }}, #66{{ onSurfaceVariant.hex }}
    inactive_colors=#ff{{ onSurface.hex }}, #ff{{ surfaceContainer.hex }}, #ff{{ surfaceContainerHigh.hex }}, #ff{{ surfaceContainerHighest.hex }}, #ff{{ surfaceDim.hex }}, #ff{{ surfaceContainerHigh.hex }}, #ff{{ onSurface.hex }}, #ff{{ onPrimary.hex }}, #ff{{ onSurface.hex }}, #ff{{ surface.hex }}, #ff{{ surfaceContainerLow.hex }}, #ff{{ shadow.hex }}, #ff{{ primary.hex }}, #ff{{ onPrimary.hex }}, #ff{{ secondary.hex }}, #ff{{ tertiary.hex }}, #ff{{ surfaceContainer.hex }}, #ff{{ onSurface.hex }}, #ff{{ inverseSurface.hex }}, #ff{{ inverseOnSurface.hex }}, #99{{ onSurfaceVariant.hex }}
  '';

  xdg.configFile."caelestia/templates/qt6ct-portal.qss".text = ''
    /*
      Hyprland share picker only.
      The binary is Qt Widgets: QMainWindow, QTabWidget, QScrollArea, QCheckBox, QPushButton.
    */

    QMainWindow,
    QDialog,
    QWidget {
      background-color: #{{ surface.hex }};
      color: #{{ onSurface.hex }};
      selection-background-color: #{{ primary.hex }};
      selection-color: #{{ onPrimary.hex }};
    }

    QLabel {
      color: #{{ onSurface.hex }};
      background: transparent;
    }

    QLabel:disabled {
      color: #{{ onSurfaceVariant.hex }};
    }

    QTabWidget::pane {
      background-color: #{{ surfaceContainerLowest.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 8px;
      top: -1px;
    }

    QTabBar::tab {
      min-height: 28px;
      padding: 6px 14px;
      margin-right: 4px;
      background-color: #{{ surfaceContainerLow.hex }};
      color: #{{ onSurfaceVariant.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-bottom-color: #{{ outlineVariant.hex }};
      border-top-left-radius: 8px;
      border-top-right-radius: 8px;
    }

    QTabBar::tab:hover {
      color: #{{ onSurface.hex }};
      background-color: #{{ surfaceContainer.hex }};
      border-color: #{{ secondary.hex }};
    }

    QTabBar::tab:selected {
      color: #{{ onSurface.hex }};
      background-color: #{{ surfaceContainerHigh.hex }};
      border-color: #{{ primary.hex }};
      border-bottom-color: #{{ surfaceContainerHigh.hex }};
    }

    QScrollArea,
    QScrollArea > QWidget,
    QFrame {
      background-color: #{{ surfaceContainerLowest.hex }};
      border: none;
    }

    QPushButton {
      min-height: 32px;
      padding: 6px 12px;
      background-color: #{{ surfaceContainer.hex }};
      color: #{{ onSurface.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 8px;
    }

    QPushButton:hover {
      background-color: #{{ surfaceContainerHigh.hex }};
      border-color: #{{ secondary.hex }};
    }

    QPushButton:pressed,
    QPushButton:checked {
      background-color: #{{ primary.hex }};
      color: #{{ onPrimary.hex }};
      border-color: #{{ primary.hex }};
    }

    QPushButton:focus {
      border-color: #{{ primary.hex }};
    }

    QCheckBox,
    QRadioButton {
      color: #{{ onSurface.hex }};
      spacing: 8px;
      background: transparent;
    }

    QCheckBox::indicator,
    QRadioButton::indicator {
      width: 16px;
      height: 16px;
      border: 1px solid #{{ outline.hex }};
      background-color: #{{ surfaceContainerLow.hex }};
    }

    QCheckBox::indicator {
      border-radius: 4px;
    }

    QRadioButton::indicator {
      border-radius: 9px;
    }

    QCheckBox::indicator:checked,
    QRadioButton::indicator:checked {
      background-color: #{{ primary.hex }};
      border-color: #{{ primary.hex }};
    }

    QScrollBar:vertical,
    QScrollBar:horizontal {
      background: transparent;
      border: none;
      margin: 0;
    }

    QScrollBar::handle:vertical,
    QScrollBar::handle:horizontal {
      background-color: #{{ outlineVariant.hex }};
      border-radius: 5px;
      min-height: 28px;
      min-width: 28px;
    }

    QScrollBar::handle:vertical:hover,
    QScrollBar::handle:horizontal:hover {
      background-color: #{{ secondary.hex }};
    }

    QToolTip {
      background-color: #{{ inverseSurface.hex }};
      color: #{{ inverseOnSurface.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 6px;
      padding: 6px;
    }
  '';

  xdg.configFile."caelestia/templates/qt6ct-caelestia.qss".text = ''
    QWidget {
      background-color: #{{ surface.hex }};
      color: #{{ onSurface.hex }};
      selection-background-color: #{{ primary.hex }};
      selection-color: #{{ onPrimary.hex }};
    }

    QDialog,
    QFileDialog,
    QMessageBox {
      background-color: #{{ surface.hex }};
    }

    QFrame,
    QGroupBox,
    QStackedWidget,
    QScrollArea,
    QMdiArea {
      background-color: #{{ surface.hex }};
      border: 0;
    }

    QMenuBar,
    QToolBar,
    QStatusBar {
      background-color: #{{ surfaceContainerLow.hex }};
      color: #{{ onSurface.hex }};
      border: 0;
      border-bottom: 1px solid #{{ outlineVariant.hex }};
      spacing: 6px;
    }

    QMenu {
      background-color: #{{ surfaceContainer.hex }};
      color: #{{ onSurface.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 8px;
      padding: 6px;
    }

    QMenu::item {
      min-height: 26px;
      padding: 6px 28px 6px 10px;
      border-radius: 6px;
      background-color: transparent;
    }

    QMenu::item:selected {
      background-color: #{{ primary.hex }};
      color: #{{ onPrimary.hex }};
    }

    QPushButton,
    QToolButton,
    QComboBox,
    QAbstractSpinBox {
      min-height: 30px;
      padding: 5px 10px;
      background-color: #{{ surfaceContainerHigh.hex }};
      color: #{{ onSurface.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 7px;
    }

    QPushButton:hover,
    QToolButton:hover,
    QComboBox:hover,
    QAbstractSpinBox:hover {
      background-color: #{{ surfaceContainerHighest.hex }};
      border-color: #{{ primary.hex }};
    }

    QPushButton:pressed,
    QPushButton:checked,
    QToolButton:pressed,
    QToolButton:checked {
      background-color: #{{ primary.hex }};
      color: #{{ onPrimary.hex }};
      border-color: #{{ primary.hex }};
    }

    QPushButton:default {
      background-color: #{{ primary.hex }};
      color: #{{ onPrimary.hex }};
      border-color: #{{ primary.hex }};
    }

    QPushButton:disabled,
    QToolButton:disabled,
    QComboBox:disabled,
    QAbstractSpinBox:disabled {
      background-color: #{{ surfaceContainerLow.hex }};
      color: #{{ onSurfaceVariant.hex }};
      border-color: #{{ outlineVariant.hex }};
    }

    QLineEdit,
    QTextEdit,
    QPlainTextEdit,
    QSpinBox,
    QDoubleSpinBox,
    QDateEdit,
    QTimeEdit,
    QDateTimeEdit {
      min-height: 30px;
      padding: 5px 9px;
      background-color: #{{ surfaceContainerLow.hex }};
      color: #{{ onSurface.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 7px;
      selection-background-color: #{{ primary.hex }};
      selection-color: #{{ onPrimary.hex }};
    }

    QLineEdit:focus,
    QTextEdit:focus,
    QPlainTextEdit:focus,
    QAbstractSpinBox:focus,
    QComboBox:focus {
      border: 1px solid #{{ primary.hex }};
      background-color: #{{ surfaceContainer.hex }};
    }

    QListView,
    QTreeView,
    QTableView,
    QColumnView {
      background-color: #{{ surface.hex }};
      alternate-background-color: #{{ surfaceContainerLow.hex }};
      color: #{{ onSurface.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 8px;
      padding: 4px;
      outline: 0;
      selection-background-color: #{{ primary.hex }};
      selection-color: #{{ onPrimary.hex }};
    }

    QListView::item,
    QTreeView::item,
    QTableView::item,
    QColumnView::item {
      min-height: 28px;
      padding: 5px 8px;
      border-radius: 6px;
    }

    QListView::item:hover,
    QTreeView::item:hover,
    QTableView::item:hover,
    QColumnView::item:hover {
      background-color: #{{ surfaceContainerHigh.hex }};
      color: #{{ onSurface.hex }};
    }

    QListView::item:selected,
    QTreeView::item:selected,
    QTableView::item:selected,
    QColumnView::item:selected {
      background-color: #{{ primary.hex }};
      color: #{{ onPrimary.hex }};
    }

    QHeaderView::section {
      min-height: 28px;
      padding: 5px 8px;
      background-color: #{{ surfaceContainer.hex }};
      color: #{{ onSurface.hex }};
      border: 0;
      border-bottom: 1px solid #{{ outlineVariant.hex }};
      border-right: 1px solid #{{ outlineVariant.hex }};
    }

    QTabWidget::pane {
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 8px;
      background-color: #{{ surface.hex }};
      top: -1px;
    }

    QTabBar::tab {
      min-height: 28px;
      padding: 6px 12px;
      color: #{{ onSurfaceVariant.hex }};
      background-color: #{{ surfaceContainerLow.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-bottom: 0;
      border-top-left-radius: 7px;
      border-top-right-radius: 7px;
      margin-right: 3px;
    }

    QTabBar::tab:selected {
      color: #{{ onSurface.hex }};
      background-color: #{{ surfaceContainerHigh.hex }};
      border-color: #{{ primary.hex }};
    }

    QTabBar::tab:hover {
      color: #{{ onSurface.hex }};
      background-color: #{{ surfaceContainerHighest.hex }};
    }

    QScrollBar {
      background-color: transparent;
      border: 0;
      margin: 2px;
    }

    QScrollBar:vertical {
      width: 10px;
    }

    QScrollBar:horizontal {
      height: 10px;
    }

    QScrollBar::handle {
      background-color: #{{ outlineVariant.hex }};
      border-radius: 5px;
      min-height: 28px;
      min-width: 28px;
    }

    QScrollBar::handle:hover {
      background-color: #{{ primary.hex }};
    }

    QScrollBar::add-line,
    QScrollBar::sub-line,
    QScrollBar::add-page,
    QScrollBar::sub-page {
      background: transparent;
      border: 0;
    }

    QSlider::groove:horizontal {
      height: 6px;
      border-radius: 3px;
      background-color: #{{ surfaceContainerHighest.hex }};
    }

    QSlider::handle:horizontal {
      width: 16px;
      margin: -5px 0;
      border-radius: 8px;
      background-color: #{{ primary.hex }};
    }

    QProgressBar {
      min-height: 8px;
      border: 0;
      border-radius: 4px;
      background-color: #{{ surfaceContainerHighest.hex }};
      text-align: center;
      color: #{{ onSurface.hex }};
    }

    QProgressBar::chunk {
      border-radius: 4px;
      background-color: #{{ primary.hex }};
    }

    QToolTip {
      color: #{{ onSurface.hex }};
      background-color: #{{ surfaceContainerHigh.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 7px;
      padding: 6px;
    }

    QWidget {
      selection-background-color: #{{ primary.hex }};
      selection-color: #{{ onPrimary.hex }};
    }

    QDialog,
    QFileDialog,
    QMessageBox {
      background-color: #{{ surface.hex }};
    }

    QMenuBar,
    QToolBar,
    QStatusBar {
      background-color: #{{ surfaceContainerLow.hex }};
      border-top: 2px solid #{{ primary.hex }};
      border-bottom: 1px solid #{{ tertiary.hex }};
    }

    QToolBar QToolButton {
      background-color: #{{ surfaceContainer.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 8px;
      padding: 5px 9px;
    }

    QToolBar QToolButton:hover {
      background-color: #{{ tertiary.hex }};
      color: #{{ onTertiary.hex }};
      border-color: #{{ tertiary.hex }};
    }

    QPushButton,
    QToolButton,
    QComboBox,
    QAbstractSpinBox {
      background-color: #{{ surfaceContainerHigh.hex }};
      border: 1px solid #{{ secondary.hex }};
      border-radius: 8px;
      padding: 6px 11px;
    }

    QPushButton:hover,
    QToolButton:hover,
    QComboBox:hover,
    QAbstractSpinBox:hover {
      background-color: #{{ tertiary.hex }};
      color: #{{ onTertiary.hex }};
      border-color: #{{ tertiary.hex }};
    }

    QPushButton:pressed,
    QPushButton:checked,
    QPushButton:default,
    QToolButton:pressed,
    QToolButton:checked {
      background-color: #{{ primary.hex }};
      color: #{{ onPrimary.hex }};
      border-color: #{{ primary.hex }};
    }

    QLineEdit,
    QTextEdit,
    QPlainTextEdit,
    QAbstractSpinBox {
      background-color: #{{ surfaceContainerLow.hex }};
      border: 1px solid #{{ secondary.hex }};
      border-radius: 8px;
    }

    QLineEdit:focus,
    QTextEdit:focus,
    QPlainTextEdit:focus,
    QAbstractSpinBox:focus,
    QComboBox:focus {
      background-color: #{{ surfaceContainer.hex }};
      border: 2px solid #{{ primary.hex }};
    }

    QListView,
    QTreeView,
    QTableView,
    QColumnView {
      background-color: #{{ surface.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-radius: 10px;
      padding: 6px;
    }

    QListView::item,
    QTreeView::item,
    QTableView::item,
    QColumnView::item {
      min-height: 30px;
      padding: 6px 9px;
      border-radius: 8px;
      border-left: 3px solid transparent;
    }

    QListView::item:hover,
    QTreeView::item:hover,
    QTableView::item:hover,
    QColumnView::item:hover {
      background-color: #{{ surfaceContainerHigh.hex }};
      border-left-color: #{{ tertiary.hex }};
    }

    QListView::item:selected,
    QTreeView::item:selected,
    QTableView::item:selected,
    QColumnView::item:selected {
      background-color: #{{ primary.hex }};
      color: #{{ onPrimary.hex }};
      border-left-color: #{{ tertiary.hex }};
    }

    QHeaderView::section {
      background-color: #{{ surfaceContainer.hex }};
      color: #{{ onSurface.hex }};
      border-bottom: 2px solid #{{ secondary.hex }};
    }

    QTabBar::tab {
      background-color: #{{ surfaceContainerLow.hex }};
      border: 1px solid #{{ outlineVariant.hex }};
      border-bottom: 2px solid #{{ secondary.hex }};
    }

    QTabBar::tab:hover {
      background-color: #{{ surfaceContainerHigh.hex }};
      border-bottom-color: #{{ tertiary.hex }};
    }

    QTabBar::tab:selected {
      background-color: #{{ surfaceContainerHigh.hex }};
      color: #{{ onSurface.hex }};
      border-color: #{{ primary.hex }};
      border-bottom-color: #{{ primary.hex }};
    }

    QScrollBar::handle {
      background-color: #{{ secondary.hex }};
      border-radius: 5px;
    }

    QScrollBar::handle:hover {
      background-color: #{{ tertiary.hex }};
    }

    QSlider::handle:horizontal,
    QSlider::handle:vertical,
    QProgressBar::chunk {
      background-color: #{{ primary.hex }};
    }
  '';
  };
}
