import 'package:flutter/material.dart';
import 'package:flutter_app/components/user_avatar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

enum FriendStatus {
  free,
  busy,
  invisible,
}

class FriendCard extends StatelessWidget {
  final String name;
  final String picture;
  final FriendStatus status;
  // curEventName is the name of the current event the friend is attending
  final String? curEventName;
  final Function(RelativeRect position) showOverflowMenu;
  final GestureTapCallback? onTap;

  final _buttonKey = GlobalKey();

  FriendCard({
    Key? key,
    required this.name,
    required this.picture,
    required this.status,
    required this.showOverflowMenu,
    required this.curEventName,
    this.onTap,
  }) : super(key: key);

  void _onButtonPressed() {
    RenderBox? box =
        _buttonKey.currentContext!.findRenderObject() as RenderBox?;
    Offset position = box!.localToGlobal(Offset.zero);
    RelativeRect rect = RelativeRect.fromLTRB(position.dx, position.dy, 0, 0);

    showOverflowMenu(rect);
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: SchejConstants.listTileDecoration,
      child: ListTile(
        onTap: onTap,
        dense: true,
        leading: UserAvatar(src: picture),
        title: Text(name, style: SchejFonts.subtitle),
        subtitle: _buildStatusText(),
        trailing: IconButton(
          key: _buttonKey,
          icon: const Icon(MdiIcons.dotsVertical),
          onPressed: _onButtonPressed,
        ),
      ),
    );
  }

  Widget _buildStatusText() {
    if (status == FriendStatus.free) {
      return Text(
        'Currently free',
        style: SchejFonts.bodyMedium.copyWith(color: SchejColors.darkGreen),
      );
    } else if (status == FriendStatus.invisible) {
      return Text(
        'Currently invisible',
        style: SchejFonts.body.copyWith(color: SchejColors.darkGray),
      );
    } else if (status == FriendStatus.busy) {
      return Text(
        'Currently at $curEventName',
        style: SchejFonts.body.copyWith(color: SchejColors.pureBlack),
      );
    }

    return const Text('');
  }
}
